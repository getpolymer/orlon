package internal

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SocketWriter struct {
	Lock *sync.Mutex
	Conn *websocket.Conn
}

func NewSocketWriter(ws *websocket.Conn) SocketWriter {
	ws.SetReadDeadline(time.Time{})
	return SocketWriter{Conn: ws, Lock: new(sync.Mutex)}
}

func (sw SocketWriter) Write(data []byte) (int, error) {
	sw.Lock.Lock()
	defer sw.Lock.Unlock()
	if err := sw.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}
