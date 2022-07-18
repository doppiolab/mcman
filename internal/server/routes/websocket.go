package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/doppiolab/mcman/internal/minecraft"
	"github.com/doppiolab/mcman/internal/minecraft/logstream"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type mcmanWsPayload struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// Serve minecraft tty service via websocket.
func ServeTerminal(mcsrv minecraft.MinecraftServer, ls logstream.LogStream) func(c echo.Context) error {
	return func(c echo.Context) error {
		ws, err := websocket.Accept(c.Response(), c.Request(), &websocket.AcceptOptions{
			// NOTE(hayeon): CompressionThreshold is set to arbitrary large value because of below issue.
			//               Maybe splitting payload can be required when there's message those lengths are longer than 16384.
			//               https://github.com/nhooyr/websocket/issues/218
			CompressionThreshold: 16384,
		})
		if err != nil {
			return errors.Wrap(err, "cannot upgrade http connection to websocket connection.")
		}
		// WebSocket conn should not be closed with below statement.
		// This is to confirm that websocket conn is closed before the http connection is ended.
		defer ws.Close(websocket.StatusInternalError, "")

		ctx := c.Request().Context()
		socketUUID := fmt.Sprintf("ws-%s", uuid.New().String())
		ls.RegisterLogCallback(socketUUID, func(lb *logstream.LogBlock) error {
			err := wsjson.Write(ctx, ws, &mcmanWsPayload{Msg: lb.Msg, Type: lb.ChanId})
			if err != nil {
				return errors.Wrap(err, "cannot send messages")
			}
			return nil
		})
		defer ls.DeregisterLogCallback(socketUUID)

		// TODO(hayeon): make the rate and duration configurable
		rateLimiter := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
		for {
			err := rateLimiter.Wait(ctx)
			if err != nil {
				return errors.Wrap(err, "rate limit error")
			}

			var payload mcmanWsPayload
			if err = wsjson.Read(ctx, ws, &payload); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}

				closeStatus := websocket.CloseStatus(err)
				if closeStatus == websocket.StatusNormalClosure || closeStatus == websocket.StatusGoingAway {
					log.Info().Msg("connection closed normaly")
					// it is okay to close the connection in this case
					return nil
				}

				// NOTE(hayeon): If payload size is longer than compression threshold.
				//               We don't want to handle this error, so disconnect the connection for now.
				if closeStatus == websocket.StatusProtocolError {
					log.Error().Err(err).Msg("cannot read messages")
					return nil
				}

				log.Error().Err(err).Msg("cannot read messages")
				continue
			}

			mcsrv.PutCommand(payload.Msg)
		}
	}
}
