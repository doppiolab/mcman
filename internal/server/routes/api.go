package routes

import (
	"net/http"

	"github.com/doppiolab/mcman/internal/minecraft/mcdata"
	"github.com/doppiolab/mcman/internal/minecraft/mcmap"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type getMapDataPayload struct {
	X int `json:"x"`
	Z int `json:"z"`
}

// Get Region List
func GetRegionList(mcDataPath string) func(c echo.Context) error {
	return func(c echo.Context) error {
		l, err := mcdata.GetRegionList(mcDataPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get region list").Error())
		}

		return c.JSON(http.StatusOK, l)
	}
}

// Get Map data for viewer
func GetMapChunkImage(mcDataPath string, tempPath string) func(c echo.Context) error {
	return func(c echo.Context) error {
		p := getMapDataPayload{}
		err := echo.PathParamsBinder(c).Int("x", &p.X).Int("z", &p.Z).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// TODO(jeongukjae): invalidate cache when the save file is changed
		var img []byte
		img, err = mcmap.MaybeCached(tempPath, p.X, p.Z)

		if err != nil {
			region, err := mcdata.GetRegion(mcDataPath, p.X, p.Z)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get level").Error())
			}

			img, err = mcmap.DrawMap(region)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to draw map").Error())
			}

			err = mcmap.Cache(tempPath, p.X, p.Z, img)
			if err != nil {
				log.Error().Err(err).Int("x", p.X).Int("z", p.Z).Msg("cannot cache")
			}
		}

		return c.Blob(http.StatusOK, "image/x-png", img)
	}
}

// Get Player data for viewer
func GetPlayerData(mcDataPath string) func(c echo.Context) error {
	return func(c echo.Context) error {
		playerData, err := mcdata.GetPlayerData(mcDataPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get player data").Error())
		}
		return c.JSON(http.StatusOK, playerData)
	}
}
