package routes

import (
	"fmt"
	"net/http"

	"github.com/doppiolab/mcman/internal/server/auth"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GetIndexPage(version string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", echo.Map{
			"Version": fmt.Sprint("#", version),
		})
	}
}

func GetLoginPage(version string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", echo.Map{
			"Version": fmt.Sprint("#", version),
		})
	}
}

func PostLoginPage(version string) func(c echo.Context) error {
	return func(c echo.Context) error {
		// TODO(hayeon): support redirect_url
		username := c.FormValue("id")
		password := c.FormValue("password")

		token, err := auth.CreateNewToken(username, password)
		if err != nil {
			return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
				"err":     err.Error(),
				"Version": fmt.Sprint("#", version),
			})
		}

		log.Info().Msg(token)

		c.SetCookie(&http.Cookie{
			Name:  auth.CookieAuthTokenKey,
			Value: token,
			Path:  "/",
		})
		return c.Redirect(http.StatusFound, "/")
	}
}
