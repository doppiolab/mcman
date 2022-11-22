package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GetModListPage(mcDataPath, version string) func(c echo.Context) error {
	modDirectory := path.Join(mcDataPath, "mods")
	return func(c echo.Context) error {
		modFilenames := []string{}

		if dirInfo, err := os.Stat(modDirectory); err == nil && dirInfo.IsDir() {
			modFiles, err := ioutil.ReadDir(modDirectory)
			if err != nil {
				log.Error().Err(err).Msg("failed to read mods directory")
				return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
					"error":   "Error reading mods directory",
					"Version": version,
				})
			}

			for _, modFile := range modFiles {
				modFilenames = append(modFilenames, modFile.Name())
			}
		}

		return c.Render(http.StatusOK, "mods.html", echo.Map{
			"Version": fmt.Sprint("#", version),
			"mods":    modFilenames,
		})
	}
}

func GetModDownloadPage(mcDataPath string) func(c echo.Context) error {
	return func(c echo.Context) error {
		var filename string
		err := echo.PathParamsBinder(c).String("filename", &filename).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if filename == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "filename is required")
		}

		modFile := path.Join(mcDataPath, "mods", filename)
		if stat, err := os.Stat(modFile); err != nil || stat.IsDir() {
			log.Error().Err(err).Str("path", modFile).Msg("failed to read mod file")
			return echo.NewHTTPError(http.StatusNotFound, "file not found")
		}

		return c.Attachment(modFile, filename)
	}
}
