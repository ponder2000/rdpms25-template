package socket

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/ponder2000/rdpms25-template/pkg/util/generic"
)

type Hub[T any, R any] struct {
	mutex sync.RWMutex

	clients   map[*Client[T, R]]bool
	broadcast chan T

	register   chan *Client[T, R]
	unregister chan *Client[T, R]
}

func NewHub[T any, R any]() *Hub[T, R] {
	h := &Hub[T, R]{
		mutex:      sync.RWMutex{},
		clients:    make(map[*Client[T, R]]bool),
		broadcast:  make(chan T),
		register:   make(chan *Client[T, R]),
		unregister: make(chan *Client[T, R]),
	}

	go h.run()
	return h
}

func (h *Hub[T, R]) Broadcast(obj T) {
	h.broadcast <- obj
}

func (h *Hub[T, R]) DisconnectClient(whereMatch func(R) bool, statusCode int, reason string) {
	h.mutex.RLock()
	vals := generic.Where(generic.GetKeysFromMap(h.clients), func(c *Client[T, R]) bool { return whereMatch(c.clientIdentifier) })
	h.mutex.RUnlock()

	if len(vals) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(vals))
	for _, v := range vals {
		go func(client *Client[T, R]) {
			defer wg.Done()
			client.closeConnection(statusCode, reason)
		}(v)
	}
	wg.Wait()
}

func (h *Hub[T, R]) run() {
	for {
		select {
		case c := <-h.register:
			h.mutex.Lock()
			h.clients[c] = true
			h.mutex.Unlock()
			slog.Info("[newConn]", "active_conn", len(h.clients))

		case c := <-h.unregister:
			h.mutex.Lock()
			delete(h.clients, c)
			h.mutex.Unlock()
			slog.Info("[disconnect]", "active_conn", len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for c := range h.clients {
				c.sendMessage(message)
			}
			h.mutex.RUnlock()

		}
	}
}

func (h *Hub[T, R]) ConnectedClients() []*Client[T, R] {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return generic.GetKeysFromMap(h.clients)
}

func (h *Hub[T, R]) Serve(
	w http.ResponseWriter, r *http.Request,
	clientBufferSize int,
	identifier R, isMsgApplicable func(T, R) bool,
	callbackClientMsg func(R, []byte), callbackConnect, callbackDisconnect func(R),
) {
	started := time.Now()
	conn, e := socketUpgrade.Upgrade(w, r, nil)
	if e != nil {
		slog.Error("unable to Serve", "err", e.Error())
		return
	}

	c := newClient[T, R](conn, h, clientBufferSize, identifier, isMsgApplicable, callbackClientMsg, callbackConnect, callbackDisconnect)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.writePump(ctx)
	go c.readPump(ctx)

	c.hub.register <- c
	c.onConnect()

	<-c.done
	c.onDisconnect()
	slog.Info("Socket connection closed..", "dur", time.Since(started))
}
