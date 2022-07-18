package world

import (
	"compress/gzip"
	"fmt"
	"os"
	"path"

	"github.com/Tnze/go-mc/save"
	"github.com/Tnze/go-mc/save/region"
	"github.com/doppiolab/mcman/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type WorldReader interface {
	// Read world/level.dat file
	GetLevel() (*save.Level, error)
	// Read world/region/r.[x].[z].mca file
	GetChunk(x, z int) (*region.Region, error)
}

type worldReader struct {
	cfg *config.MinecraftConfig
}

func NewReader(cfg *config.MinecraftConfig) WorldReader {
	return &worldReader{
		cfg: cfg,
	}
}

func (r *worldReader) GetLevel() (*save.Level, error) {
	levelDatPath := path.Join(r.cfg.WorkingDir, "world", "level.dat")
	f, err := os.Open(levelDatPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find level.dat file")
	}

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

func (wr *worldReader) GetChunk(x, z int) (*region.Region, error) {
	chunkFilePath := path.Join(wr.cfg.WorkingDir, "world", "region", fmt.Sprintf("r.%d.%d.mca", x, z))
	r, err := region.Open(chunkFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open chunk file")
	}
	defer r.Close()

	var c save.Chunk
	data, err := r.ReadSector(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open sector")
	}
	err = c.Load(data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load sector data")
	}
	// for i := len(c.Sections) - 1 ; i >= 0;i-- {
	// 	log.Info().Interface("section", c.Sections[i]).Msg("H")
	// }
	log.Info().Interface("c", c).Msg("Hello")
	return r, nil
}
