package notification

import (
	"log"
	"net/http"
	"sync"

	"github.com/ailisp/reallyfastci/core"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]chan bool
}{
	connections: make(map[*websocket.Conn]chan bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }}

func NotifyWebSocket(event *core.BuildEvent) {
	connectionPool.RLock()
	defer connectionPool.RUnlock()
	for ws := range connectionPool.connections {
		if err := ws.WriteJSON(map[string]string{
			"commit": event.Commit,
			"status": core.BuildStatusStr(event.Status),
		}); err != nil {
			log.Printf("Notify error: %v", err)
			connectionPool.Lock()
			// connectionPool.connections[ws] <- true
			delete(connectionPool.connections, ws)
			connectionPool.Unlock()
		}
	}
}

func WebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	exitChan := make(chan bool)
	connectionPool.Lock()
	connectionPool.connections[ws] = exitChan
	connectionPool.Unlock()

	defer func(connection *websocket.Conn) {
		connectionPool.Lock()
		delete(connectionPool.connections, connection)
		connectionPool.Unlock()
	}(ws)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
	}

	// <-exitChan
	return nil
}
