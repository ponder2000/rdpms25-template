package socket

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type Client[T any, R any] struct {
	conn *websocket.Conn
	hub  *Hub[T, R]

	send chan T
	done chan bool

	clientIdentifier R
	isMsgApplicable  func(obj T, identifier R) bool

	receiveMsgCallback func(identifier R, data []byte)
	disconnectCallback func(identifier R)
	connectCallback    func(identifier R)
}

func newClient[T any, R any](
	conn *websocket.Conn, hub *Hub[T, R],
	bufferSize int,
	identifier R, isMsgApplicable func(T, R) bool,
	callbackClientMsg func(R, []byte), callbackConnect, callbackDisconnect func(R),
) *Client[T, R] {
	c := Client[T, R]{
		conn: conn,
		hub:  hub,

		send: make(chan T, bufferSize),
		done: make(chan bool),

		clientIdentifier: identifier,
		isMsgApplicable:  isMsgApplicable,

		connectCallback:    callbackConnect,
		disconnectCallback: callbackDisconnect,
		receiveMsgCallback: callbackClientMsg,
	}
	return &c
}

func (c *Client[T, R]) sendMessage(obj T) {
	if c.isMsgApplicable != nil {
		if c.isMsgApplicable(obj, c.clientIdentifier) {
			select {
			case c.send <- obj:
			default:
				slog.Error("unable to write message to client. Buffer overflow", "identifer", c.clientIdentifier)
			}
		}
	} else {
		// in case of no filter added then all events will be sent
		c.send <- obj
	}
}

func (c *Client[T, R]) receivedMsgFromClient(data []byte) {
	slog.Debug("Msg from Client callback", "data", string(data))
	if c.receiveMsgCallback != nil {
		go c.receiveMsgCallback(c.clientIdentifier, data)
	}
}

func (c *Client[T, R]) onDisconnect() {
	slog.Debug("Client disconnected callback", "addr", c.conn.RemoteAddr())
	if c.disconnectCallback != nil {
		c.disconnectCallback(c.clientIdentifier)
	}
}

func (c *Client[T, R]) onConnect() {
	slog.Debug("Client connected callback", "addr", c.conn.RemoteAddr())
	if c.connectCallback != nil {
		c.connectCallback(c.clientIdentifier)
	}
}

func (c *Client[T, R]) closeConnection(statusCode int, reason string) {
	c.hub.unregister <- c

	closeMessage := websocket.FormatCloseMessage(statusCode, reason)
	e := c.conn.WriteMessage(websocket.CloseMessage, closeMessage)
	if e != nil {
		slog.Error("Close message error to client", "err", e.Error())
	}
	_ = c.conn.Close()
	c.done <- true
}

func (c *Client[T, R]) readPump(ctx context.Context) {
	defer func() {
		slog.Info("Closing read channel")
		c.hub.unregister <- c
		_ = c.conn.Close()
		c.done <- true
	}()
	_ = c.conn.SetReadDeadline(time.Now().Add(PongWait))

	c.conn.SetPongHandler(func(string) error {
		slog.Debug("Pong from Client")
		_ = c.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, e := c.conn.ReadMessage()
			if e != nil {
				if websocket.IsUnexpectedCloseError(e, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseProtocolError) {
					slog.Error("Unexpected close error in read channel", "err", e.Error())
				}
				return
			}
			c.receivedMsgFromClient(message)
		}
	}
}

func (c *Client[T, R]) writePump(ctx context.Context) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		slog.Info("Closing write channel")
		ticker.Stop()
		_ = c.conn.Close()
		c.done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if e := c.conn.WriteJSON(message); e != nil {
				slog.Error("error writing to client", "err", e)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if e := c.conn.WriteMessage(websocket.PingMessage, nil); e != nil {
				return
			}
		}
	}
}
