package mcdata

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Tnze/go-mc/save"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	UUID string
	Name string
	X    float64
	Y    float64
	Z    float64
}

// GetPlayerData returns the all player data of the given datapath.
//
// This function uses the ${dataPath}/world/playerdata/*.dat files and
// ${dataPath}/usercache.json file.
func GetPlayerData(dataPath string) ([]PlayerData, error) {
	playerdataPath := path.Join(dataPath, "world", "playerdata")
	files, err := ioutil.ReadDir(playerdataPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list directory")
	}

	userNameMap, err := readUsernameMapFromUserCache(path.Join(dataPath, "usercache.json"))
	if err != nil {
		return nil, errors.Wrap(err, "cannot read usercache.json")
	}

	results := []PlayerData{}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".dat") {
			continue
		}

		filepath := path.Join(playerdataPath, file.Name())
		data, err := readPlayerData(filepath)
		if err != nil {
			log.Error().Err(err).Str("file", file.Name()).Msg("cannot read playerdata")
			continue
		}

		uuid := strings.TrimSuffix(file.Name(), ".dat")
		name, ok := userNameMap[uuid]
		if !ok {
			name = "**UNKNOWN**"
		}

		results = append(results, PlayerData{
			UUID: uuid,
			Name: name,
			X:    data.Pos[0],
			Y:    data.Pos[1],
			Z:    data.Pos[2],
		})
	}

	return results, nil
}

// Reads the usercache.json file and returns a uuid to username map.
func readUsernameMapFromUserCache(userCacheFile string) (map[string]string, error) {
	f, err := os.Open(userCacheFile)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open user cache")
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read user cache")
	}

	var users []struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	}
	if err := json.Unmarshal(byteValue, &users); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal user cache")
	}

	userNameMap := map[string]string{}
	for _, user := range users {
		userNameMap[user.UUID] = user.Name
	}

	return userNameMap, nil
}

// Reads playerdata from a gzipped file.
func readPlayerData(filepath string) (*save.PlayerData, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open playerdata file")
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read playerdata file as gzip")
	}
	defer r.Close()

	data, err := save.ReadPlayerData(r)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal playerdata file")
	}

	return &data, nil
}
