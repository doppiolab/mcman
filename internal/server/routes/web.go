package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetIndexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
