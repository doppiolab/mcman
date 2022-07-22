package mcmap

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strings"

	"github.com/doppiolab/mcman/internal/minecraft/mcdata"
	"github.com/pkg/errors"
)

// Draw region info as map image in PNG format
func DrawMap(r *mcdata.RegionInfo) ([]byte, error) {
	topLeft := image.Point{0, 0}
	numRowPerRegion := mcdata.NumChunkRows * mcdata.NumBlockRows
	bottomRight := image.Point{numRowPerRegion, numRowPerRegion}

	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	originX := r.RegionX * numRowPerRegion
	originZ := r.RegionZ * numRowPerRegion
	notFoundBlocks := map[string]bool{}

	for _, tvc := range r.Chunks {
		chunkX := int(tvc.X*mcdata.NumBlockRows) - originX
		chunkZ := int(tvc.Z*mcdata.NumBlockRows) - originZ
		for z, blocks := range tvc.Blocks {
			for x, block := range blocks {
				color, ok := colorMap[block.ID]
				if !ok {
					notFoundBlocks[block.ID] = true
					continue
				}

				posX := chunkX + x
				posZ := chunkZ + z

				targetColor := color
				ratio := (float32(block.SkyLight+block.BlockLight)/32.0)*0.5 + 0.5
				targetColor.R = uint8(float32(targetColor.R) * ratio)
				targetColor.G = uint8(float32(targetColor.G) * ratio)
				targetColor.B = uint8(float32(targetColor.B) * ratio)
				img.Set(posX, posZ, targetColor)
			}
		}
	}

	if len(notFoundBlocks) != 0 {
		keys := make([]string, 0, len(notFoundBlocks))
		for k := range notFoundBlocks {
			keys = append(keys, k)
		}
		return nil, errors.New(fmt.Sprintf("cannot find color for %s", strings.Join(keys, ", ")))
	}

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode image")
	}

	return buf.Bytes(), nil
}
