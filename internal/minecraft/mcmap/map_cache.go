package mcmap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

// Return cached image if exists, otherwise return error
func MaybeCached(tempDirectory string, x, z int) ([]byte, error) {
	fileInfo, err := os.Stat(path.Join(tempDirectory, imageCacheName(x, z)))
	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.Wrap(err, "cannot find cache file")
	}

	// TODO(jeongukjae): make this configurable
	if time.Since(fileInfo.ModTime()) > time.Minute*15 {
		return nil, errors.New("cache file is too old")
	}

	file, err := os.Open(path.Join(tempDirectory, imageCacheName(x, z)))
	if err != nil {
		return nil, errors.Wrap(err, "cannot open cache file")
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

// Cache image to temp directory
func Cache(tempDirectory string, x, z int, img []byte) error {
	file, err := os.Create(path.Join(tempDirectory, imageCacheName(x, z)))
	if err != nil {
		return errors.Wrap(err, "cannot create cache file")
	}
	defer file.Close()

	_, err = file.Write(img)
	if err != nil {
		return errors.Wrap(err, "cannot write cache file")
	}

	return nil
}

func imageCacheName(x, z int) string {
	return fmt.Sprintf("map.%d.%d.png", x, z)
}
