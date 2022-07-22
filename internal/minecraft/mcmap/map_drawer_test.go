package mcmap

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/doppiolab/mcman/internal/minecraft/mcdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDrawMap_SnapshotTest(t *testing.T) {
	tests := []struct{ X, Z int }{
		{-1, -1},
		{-1, 0},
		{0, -1},
		{0, 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test - %d,%d", tt.X, tt.Z), func(t *testing.T) {
			regionInfo, err := mcdata.GetRegion("../testdata", tt.X, tt.Z)
			require.NoError(t, err)

			snapsnotFilename := fmt.Sprintf("./snapshots/map.%d.%d.png", tt.X, tt.Z)

			pngData, err := DrawMap(regionInfo)
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(snapsnotFilename)
			require.NoError(t, err)
			assert.Equal(t, expected, pngData)
		})
	}
}
