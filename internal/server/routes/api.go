package routes

import (
	"net/http"

	"github.com/doppiolab/mcman/internal/minecraft/world"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type getMapDataPayload struct {
	X int `json:"x"`
	Z int `json:"z"`
}

// Get Map data for viewer
func GetMapData(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		p := getMapDataPayload{}
		if err := c.Bind(&p); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		region, err := reader.GetRegion(p.X, p.Z)
		if err != nil {
			return errors.Wrap(err, "failed to get level")
		}
		return c.JSON(http.StatusOK, region)
	}
}

// Get Player data for viewer
func GetPlayerData(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		return nil
	}
}
