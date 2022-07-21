package world

import (
	"compress/gzip"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/Tnze/go-mc/level"
	"github.com/Tnze/go-mc/level/block"
	"github.com/Tnze/go-mc/save"
	"github.com/Tnze/go-mc/save/region"
	"github.com/doppiolab/mcman/internal/config"
	"github.com/pkg/errors"
)

type BlockInfo struct {
	ID         string // block id string ex> "minecraft:stone"
	BlockLight uint8  // the amount of block-emitted light
	SkyLight   uint8  // the amount of sunlight or moonlight hitting each block
}

type TopViewChunk struct {
	X      int32         // bottom left corner coordinate of region
	Z      int32         // bottom left corner coordinate of region
	Blocks [][]BlockInfo // block infos. [x][z], 16*16
}

// TopViewRegion covers all 32*32 chunks in a single region file.
type TopViewRegion struct {
	RegionX int
	RegionZ int
	Chunks  []*TopViewChunk // chunk infos.
}

type WorldReader interface {
	// Read world/level.dat file
	GetLevel() (*save.Level, error)
	// Read world/region/r.[x].[z].mca file
	GetRegion(x, z int) (*TopViewRegion, error)
}

type worldReader struct {
	cfg      *config.MinecraftConfig
	readLock sync.Mutex
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

const numValuesPerHeightmap = 16 * 16
const numBitsPerValueOfHeightmap = 9

func (wr *worldReader) GetRegion(x, z int) (*TopViewRegion, error) {
	chunkFilePath := path.Join(wr.cfg.WorkingDir, "world", "region", fmt.Sprintf("r.%d.%d.mca", x, z))
	r, err := region.Open(chunkFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open chunk file")
	}
	defer r.Close()

	result := &TopViewRegion{
		RegionX: x,
		RegionZ: z,
	}

	var wg sync.WaitGroup
	queue := make(chan *TopViewChunk, 1)
	wg.Add(32 * 32)
	for z := 0; z < 32; z++ {
		go func(z int) {
			for x := 0; x < 32; x++ {
				queue <- wr.getChunkData(r, x, z)
			}
		}(z)
	}

	go func() {
		for chunk := range queue {
			if chunk != nil {
				result.Chunks = append(result.Chunks, chunk)
			}
			wg.Done()
		}
	}()

	wg.Wait()
	close(queue)

	return result, nil
}

// Get Chunk data for viewer.
//
// NOTE(jeongukjae): this function will return nil if raise error
func (wr *worldReader) getChunkData(r *region.Region, regionX, regionZ int) *TopViewChunk {
	c, err := wr.getChunkFromRegion(r, regionX, regionZ)
	if err != nil {
		return nil
	}

	if !isValidChunk(c) {
		return nil
	}

	lc, err := level.ChunkFromSave(c)
	if err != nil {
		return nil
	}

	result := &TopViewChunk{
		X:      c.XPos,
		Z:      c.ZPos,
		Blocks: make([][]BlockInfo, 16),
	}

	for z := 0; z < 16; z++ {
		result.Blocks[z] = make([]BlockInfo, 16)
		for x := 0; x < 16; x++ {
			y := lc.HeightMaps.WorldSurface.Get(z*16+x) - 1
			section := lc.Sections[y/16]
			yPos := y % 16

			pos := yPos*16*16 + z*16 + x

			if section.BlockLight != nil {
				result.Blocks[z][x].BlockLight = (section.BlockLight[pos/2] >> ((pos % 2) * 4)) & 0x0F
			}
			if section.SkyLight != nil {
				result.Blocks[z][x].SkyLight = (section.SkyLight[pos/2] >> ((pos % 2) * 4)) & 0x0F
			}

			blockStateID := section.GetBlock(pos)
			result.Blocks[z][x].ID = block.StateList[blockStateID].ID()
		}
	}

	return result
}

// Get chunk from x, z coordinates.
func (wr *worldReader) getChunkFromRegion(r *region.Region, x, z int) (*save.Chunk, error) {
	var c save.Chunk

	if !r.ExistSector(x, z) {
		return nil, errors.New("cannot find chunk")
	}

	wr.readLock.Lock()
	data, err := r.ReadSector(x, z)
	wr.readLock.Unlock()
	if err != nil {
		return nil, errors.Wrap(err, "cannot open sector")
	}

	err = c.Load(data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load sector data")
	}

	return &c, nil
}

// Return true if chunk is valid.
func isValidChunk(c *save.Chunk) bool {
	if c.Status == "full" {
		return true
	}

	// TODO(jeongukjae): add futher checks

	return false
}
