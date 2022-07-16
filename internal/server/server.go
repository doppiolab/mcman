package server

import (
	"net/http"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/labstack/echo/v4"
)

func New(cfg *config.ServerConfig) (*http.Server, error) {
	e := echo.New()
	e.Static("/static", "static")

	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: e,
	}
	return httpServer, nil
}
