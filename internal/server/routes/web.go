package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetIndexPage() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	}
}
