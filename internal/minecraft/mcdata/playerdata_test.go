package mcdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayerData(t *testing.T) {
	playerData, err := GetPlayerData("../testdata")

	assert.NoError(t, err)
	assert.NotNil(t, playerData)
	assert.Equal(t, 1, len(playerData))
	assert.Equal(t, "test-uuid", playerData[0].UUID)
	assert.Equal(t, "jeongukjae", playerData[0].Name)
}

func TestReadUsernameMapFromUserCache_WithTestData(t *testing.T) {
	data, err := readUsernameMapFromUserCache("../testdata/usercache.json")

	name, ok := data["test-uuid"]

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.True(t, ok)
	assert.Equal(t, "jeongukjae", name)
}

func TestReadPlayerData_WithTestData(t *testing.T) {
	data, err := readPlayerData("../testdata/world/playerdata/test-uuid.dat")
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
