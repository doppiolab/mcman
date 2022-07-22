package mcmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaybeCache_InvalidPath(t *testing.T) {
	data, err := MaybeCached("./invalid", 0, 0)

	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestMaybeCache_FetchCachedData(t *testing.T) {
	temp := t.TempDir()
	data := []byte{1, 2, 3}
	err := Cache(temp, 0, 0, data)
	assert.NoError(t, err)

	cached, err := MaybeCached(temp, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, data, cached)
}
