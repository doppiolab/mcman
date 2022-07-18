package routes

import (
	"github.com/doppiolab/mcman/internal/minecraft/world"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Get Map data for viewer
func GetMapData(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		// _, err := reader.GetLevel()
		_, err := reader.GetChunk(0, -1)
		if err != nil {
			log.Error().Err(err).Msg("cannot get chunk")
			return errors.Wrap(err, "failed to get level")
		}
		return nil
	}
}

// Get Player data for viewer
func GetPlayerData(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		return nil
	}
}
