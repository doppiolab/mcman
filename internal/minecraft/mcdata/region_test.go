package mcdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegionList(t *testing.T) {
	regions, err := GetRegionList("../testdata")

	assert.NoError(t, err)
	assert.Equal(t, 4, len(regions))
}

func TestGetRegionList_WithInvalidPath(t *testing.T) {
	_, err := GetRegionList("./invalid")

	assert.Error(t, err)
}

func TestGetRegion_ReadWithoutError(t *testing.T) {
	region, err := GetRegion("../testdata", 0, 0)

	assert.NoError(t, err)
	assert.NotNil(t, region)
}

func TestGetRegion_ReadInvalidXZ(t *testing.T) {
	_, err := GetRegion("../testdata", 30, 30)

	assert.Error(t, err)
}
