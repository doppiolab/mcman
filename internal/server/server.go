package server

import (
	"html/template"
	"net/http"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/labstack/echo/v4"
)

func New(cfg *config.ServerConfig) (*http.Server, error) {
	e := echo.New()
	renderer := &templateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "static")

	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: e,
	}
	return httpServer, nil
}
