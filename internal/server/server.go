package server

import (
	"net/http"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/labstack/echo/v4"
)

func New(cfg *config.ServerConfig) (*http.Server, error) {
	echoHandler := echo.New()

	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: echoHandler,
	}
	return httpServer, nil
}
