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

// Get Region List
func GetRegionList(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		l, err := reader.GetRegionList()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get region list").Error())
		}

		return c.JSON(http.StatusOK, l)
	}
}

// Get Map data for viewer
func GetMapChunkImage(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		p := getMapDataPayload{}
		err := echo.PathParamsBinder(c).Int("x", &p.X).Int("z", &p.Z).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// TODO(jeongukjae): Cache and regenerate if changed
		region, err := reader.GetRegion(p.X, p.Z)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get level").Error())
		}

		img, err := world.DrawMap(region)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to draw map").Error())
		}

		return c.Blob(http.StatusOK, "image/x-png", img)
	}
}

// Get Player data for viewer
func GetPlayerData(reader world.WorldReader) func(c echo.Context) error {
	return func(c echo.Context) error {
		playerData, err := reader.GetPlayerData()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get player data").Error())
		}
		return c.JSON(http.StatusOK, playerData)
	}
}
