package routes

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// upgrade http connection to websocket
var httpToWsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Serve minecraft tty service via websocket.
func ServeTerminal() func(c echo.Context) error {
	return func(c echo.Context) error {
		ws, err := httpToWsUpgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return errors.Wrap(err, "cannot upgrade http connection to websocket connection.")
		}
		defer ws.Close()
	}
}
