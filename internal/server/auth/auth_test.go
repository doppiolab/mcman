package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomPassword(t *testing.T) {
	// NOTE(hayeon): It is hard to test randomness, so I just check
	//               this function can return password without any error.
	for i := 0; i < 100; i++ {
		generateRandomPassword()
	}
}

func TestMaybeGenerateRandomPassword(t *testing.T) {
	envKey := "MCMAN_INIT_PASSWORD"
	password := "test-password"

	retrievedBeforeSet := MaybeGenerateRandomPassword(envKey)
	err := os.Setenv(envKey, password)
	require.NoError(t, err)
	retrievedAfterSet := MaybeGenerateRandomPassword(envKey)

	assert.NotEqual(t, password, retrievedBeforeSet)
	assert.Equal(t, password, retrievedAfterSet)
}

func TestMaybeGenerateRandomPassword_InvalidFormat(t *testing.T) {
	envKey := "MCMAN_INIT_PASSWORD"
	password := "" // zero length string cannot be used as password

	err := os.Setenv(envKey, password)
	require.NoError(t, err)
	retrieved := MaybeGenerateRandomPassword(envKey)

	assert.NotEqual(t, password, retrieved)
}
