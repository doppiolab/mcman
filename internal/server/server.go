package server

import (
	"html/template"
	"net/http"
	"path"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/doppiolab/mcman/internal/logstream"
	"github.com/doppiolab/mcman/internal/minecraft"
	"github.com/doppiolab/mcman/internal/server/auth"
	"github.com/doppiolab/mcman/internal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func New(
	cfg *config.ServerConfig,
	mcsrv minecraft.Server,
	mcDataPath string,
	ls logstream.LogStream) (*http.Server, error) {
	auth.Initialize(cfg)

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Err(v.Error).
				Msg("request")
			return nil
		},
	}))

	renderer := &templateRenderer{
		templates: template.Must(template.ParseGlob(path.Join(cfg.TemplatePath, "*.html"))),
	}
	e.Renderer = renderer

	e.Static("/static", cfg.StaticPath)
	e.File("/favicon.ico", path.Join(cfg.StaticPath, "favicon.ico"))

	e.GET("/login", routes.GetLoginPage())
	e.POST("/login", routes.PostLoginPage())

	authMiddleware := auth.NewJWTMiddleware()

	e.GET("/", routes.GetIndexPage(), authMiddleware)
	e.GET("/ws/terminal", routes.ServeTerminal(mcsrv, ls), authMiddleware)

	e.GET("/api/v1/regions", routes.GetRegionList(mcDataPath), authMiddleware)
	e.GET("/api/v1/chunk/:x/:z/map.png", routes.GetMapChunkImage(mcDataPath, cfg.TemporaryPath), authMiddleware)
	e.GET("/api/v1/players", routes.GetPlayerData(mcDataPath), authMiddleware)

	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: e,
	}
	return httpServer, nil
}
