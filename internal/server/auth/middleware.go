package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

const CookieAuthTokenKey = "authtoken"

const (
	secretKeyCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIZKLMNOPQSTUVWXYZ0123456789-_@#!%$"
	secretKeyLength  = 64
)

var signingSecretKey []byte

func init() {
	signingSecretKey = []byte(generateRandomString(secretKeyCharset, secretKeyLength))
}

// Create new JWT Middleware.
func NewJWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    signingSecretKey,
		SigningMethod: jwt.SigningMethodHS256.Name,
		TokenLookup:   fmt.Sprintf("cookie:%s", CookieAuthTokenKey),
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != jwt.SigningMethodHS256.Name {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return signingSecretKey, nil
			}

			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, errors.Wrap(err, "cannot parse token")
			}

			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			return token, nil
		},
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		},
	})
}
