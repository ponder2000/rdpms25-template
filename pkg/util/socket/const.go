package socket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// socket const
const (
	WriteWait       = 10 * time.Second
	PongWait        = 60 * time.Second
	PingPeriod      = (PongWait * 9) / 10
	ReadBufferSize  = 1024
	WriteBufferSize = 1024 * 4
)

var NewLine = []byte{'\n'}

var socketUpgrade = websocket.Upgrader{
	ReadBufferSize:    ReadBufferSize,
	WriteBufferSize:   WriteBufferSize,
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: true,
}
