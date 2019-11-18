package notification

import (
	"fmt"

	"github.com/ailisp/reallyfastci/core"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func NotifyWebSocket(event *core.BuildEvent) {

}

func WebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			err := websocket.Message.Send(ws, "Hello, Client!")
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
