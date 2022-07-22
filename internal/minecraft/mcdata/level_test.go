package mcdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLevel_WithTestData(t *testing.T) {
	level, err := GetLevel("../testdata")

	assert.NoError(t, err)
	assert.NotNil(t, level)
}

func TestGetLevel_ShouldFailWithInvalidPath(t *testing.T) {
	_, err := GetLevel("./invalid-path")

	assert.Error(t, err)
}
