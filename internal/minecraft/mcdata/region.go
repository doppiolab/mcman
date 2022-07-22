package mcdata

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"

	"github.com/Tnze/go-mc/level"
	"github.com/Tnze/go-mc/level/block"
	"github.com/Tnze/go-mc/save"
	"github.com/Tnze/go-mc/save/region"
	"github.com/pkg/errors"
)

const (
	NumChunkRows = 32
	NumBlockRows = 16
)

type RegionPoint struct {
	X int
	Z int
}

// Get all region list in dataPath.
//
// Valid region data file should be located in "world/region" directory and have
// filename like "r.x.z.mca".
func GetRegionList(dataPath string) ([]RegionPoint, error) {
	chunkFilePath := path.Join(dataPath, "world", "region")
	files, err := ioutil.ReadDir(chunkFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list directory")
	}

	results := []RegionPoint{}

	for _, file := range files {
		var x, z int
		_, err := fmt.Sscanf(file.Name(), "r.%d.%d.mca", &x, &z)
		if err == nil {
			results = append(results, RegionPoint{x, z})
		}
	}

	return results, nil
}

// RegionInfo covers all 32*32 chunks in a single region file.
type RegionInfo struct {
	RegionX int
	RegionZ int
	Chunks  []*ChunkInfo // chunk infos.
}

type ChunkInfo struct {
	X      int           // bottom left corner coordinate of region
	Z      int           // bottom left corner coordinate of region
	Blocks [][]BlockInfo // block infos. [x][z], 16*16
}

type BlockInfo struct {
	ID         string // block id string ex> "minecraft:stone"
	BlockLight uint8  // the amount of block-emitted light
	SkyLight   uint8  // the amount of sunlight or moonlight hitting each block
}

// Read and parse region data for given x, z.
func GetRegion(dataPath string, regionX, regionZ int) (*RegionInfo, error) {
	chunkFilePath := path.Join(dataPath, "world", "region", regionFilename(regionX, regionZ))
	r, err := region.Open(chunkFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open chunk file")
	}
	defer r.Close()

	result := &RegionInfo{
		RegionX: regionX,
		RegionZ: regionZ,
	}

	queue := make(chan *ChunkInfo, 1)

	var wg sync.WaitGroup
	wg.Add(NumChunkRows * NumChunkRows)

	// Parallelize chunk parsing.
	regionData := &regionData{Region: r}
	for z := 0; z < NumChunkRows; z++ {
		go func(z int) {
			for x := 0; x < NumChunkRows; x++ {
				chunk, err := getChunkData(regionData, x, z)

				if err != nil {
					queue <- nil
				} else {
					queue <- chunk
				}
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

type regionData struct {
	Region *region.Region
	lock   sync.Mutex
}

func regionFilename(x, z int) string {
	return fmt.Sprintf("r.%d.%d.mca", x, z)
}

// Read chunk data and parse for block info.
func getChunkData(r *regionData, chunkX, chunkZ int) (*ChunkInfo, error) {
	c, err := getChunkFromRegion(r, chunkX, chunkZ)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get chunk")
	}

	if !isValidChunk(c) {
		return nil, errors.New("invalid chunk data")
	}

	lc, err := level.ChunkFromSave(c)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get chunk")
	}

	result := &ChunkInfo{
		X:      int(c.XPos),
		Z:      int(c.ZPos),
		Blocks: make([][]BlockInfo, NumBlockRows),
	}

	for z := 0; z < NumBlockRows; z++ {
		result.Blocks[z] = make([]BlockInfo, NumBlockRows)

		for x := 0; x < NumBlockRows; x++ {
			y := lc.HeightMaps.WorldSurface.Get(z*NumBlockRows+x) - 1
			section := lc.Sections[y/NumBlockRows]
			y = y % NumBlockRows
			offset := y*NumBlockRows*NumBlockRows + z*NumBlockRows + x

			if section.BlockLight != nil {
				result.Blocks[z][x].BlockLight = (section.BlockLight[offset/2] >> ((offset % 2) * 4)) & 0x0F
			}
			if section.SkyLight != nil {
				result.Blocks[z][x].SkyLight = (section.SkyLight[offset/2] >> ((offset % 2) * 4)) & 0x0F
			}

			blockStateID := section.GetBlock(offset)
			result.Blocks[z][x].ID = block.StateList[blockStateID].ID()
		}
	}

	return result, nil
}

// Wrapper funciton to read chunk from region.
func getChunkFromRegion(r *regionData, x, z int) (*save.Chunk, error) {
	var c save.Chunk

	if !r.Region.ExistSector(x, z) {
		return nil, errors.New("cannot find chunk")
	}

	// use mutex to avoid concurrent read
	r.lock.Lock()
	data, err := r.Region.ReadSector(x, z)
	r.lock.Unlock()
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

	return false
}
