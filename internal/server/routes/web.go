package routes

import (
	"fmt"
	"net/http"

	"github.com/doppiolab/mcman/internal/server/auth"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GetIndexPage(gitCommit string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", echo.Map{
			"GitCommit": fmt.Sprint("#", gitCommit),
		})
	}
}

func GetLoginPage(gitCommit string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", echo.Map{
			"GitCommit": fmt.Sprint("#", gitCommit),
		})
	}
}

func PostLoginPage(gitCommit string) func(c echo.Context) error {
	return func(c echo.Context) error {
		// TODO(hayeon): support redirect_url
		username := c.FormValue("id")
		password := c.FormValue("password")

		token, err := auth.CreateNewToken(username, password)
		if err != nil {
			return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
				"err":       err.Error(),
				"GitCommit": fmt.Sprint("#", gitCommit),
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
