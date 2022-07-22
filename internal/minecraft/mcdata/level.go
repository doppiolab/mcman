package mcdata

import (
	"compress/gzip"
	"os"
	"path"

	"github.com/Tnze/go-mc/save"
	"github.com/pkg/errors"
)

// GetLevel returns the level data of the given datapath.
//
// This function uses the ${dataPath}/world/level.dat file.
func GetLevel(dataPath string) (*save.Level, error) {
	levelDatPath := path.Join(dataPath, "world", "level.dat")
	f, err := os.Open(levelDatPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find level.dat file")
	}
	defer f.Close()

	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read level.dat file as gzip")
	}
	defer gzipReader.Close()

	data, err := save.ReadLevel(gzipReader)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal level.dat file")
	}

	return &data, nil
}
