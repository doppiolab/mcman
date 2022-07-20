package server

import (
	"html/template"
	"net/http"
	"path"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/doppiolab/mcman/internal/minecraft"
	"github.com/doppiolab/mcman/internal/minecraft/logstream"
	"github.com/doppiolab/mcman/internal/minecraft/world"
	"github.com/doppiolab/mcman/internal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func New(
	cfg *config.ServerConfig,
	mcsrv minecraft.MinecraftServer,
	ls logstream.LogStream,
	worldReader world.WorldReader) (*http.Server, error) {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().Str("method", v.Method).Str("uri", v.URI).Int("status", v.Status).Dur("latency", v.Latency).Msg("request")
			return nil
		},
	}))

	renderer := &templateRenderer{
		templates: template.Must(template.ParseGlob(path.Join(cfg.TemplatePath, "*.html"))),
	}
	e.Renderer = renderer
	e.Static("/static", cfg.StaticPath)
	e.File("/favicon.ico", path.Join(cfg.StaticPath, "favicon.ico"))

	e.GET("/", routes.GetIndexPage())
	e.GET("/ws/terminal", routes.ServeTerminal(mcsrv, ls))

	e.POST("/api/v1/map", routes.GetMapData(worldReader))
	e.POST("/api/v1/player", routes.GetPlayerData(worldReader))

	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: e,
	}
	return httpServer, nil
}
