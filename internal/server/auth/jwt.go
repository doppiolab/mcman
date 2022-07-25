package auth

import (
	"time"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	userID   string
	password string
)

func Initialize(cfg *config.ServerConfig) {
	// TODO(hayeon): implement multi-user feature
	userID = "admin"
	password = MaybeGenerateRandomPassword(cfg.PasswordEnvKey)

	log.Info().Msgf("You can login with user id: '%s' and password: '%s'", userID, password)
}

type jwtCustomClaims struct {
	jwt.StandardClaims
}

func CreateNewToken(inputID, inputPassword string) (string, error) {
	if inputID != userID || inputPassword != password {
		return "", errors.New("invalid user id or password")
	}

	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			// TODO(hayeon): make this configurable
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(signingSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "cannot sign token")
	}

	return t, nil
}
