package minecraft

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMinecraftCommand(t *testing.T) {
	tests := []struct {
		cfg      *config.MinecraftConfig
		expected []string
	}{
		{
			cfg: &config.MinecraftConfig{
				JarPath: "server.jar",
				Args:    []string{"nogui"},
			},
			expected: []string{"-jar", "server.jar", "nogui"},
		},
		{
			cfg: &config.MinecraftConfig{
				JavaOptions: []string{"-Xms1024m"},
				JarPath:     "server2.jar",
				Args:        []string{"nogui"},
			},
			expected: []string{"-Xms1024m", "-jar", "server2.jar", "nogui"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			actual := createMinecraftCommandArgs(tt.cfg)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStreamReaderToChan(t *testing.T) {
	tests := []string{
		"dummy_string",
		"안녕하세요, 한글 로그에 대해서도 잘 동작하기를 바래요~",
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			r := strings.NewReader(tt)
			outChan := make(chan string, 100)

			go streamReaderToChan(r, outChan)

			actual := ""
			for d := range outChan {
				actual += d
			}

			assert.Equal(t, tt, actual)
		})
	}
}

func TestMaybeCreateWorkingDir(t *testing.T) {
	tempDir := t.TempDir()

	// create working dir
	notExistedPath := path.Join(tempDir, "server-data")
	err := maybeCreateWorkingDir(notExistedPath)
	require.NoError(t, err)
	stat, err := os.Stat(notExistedPath)
	require.NoError(t, err)
	require.True(t, stat.IsDir())

	// Even if maybeCreateWorkingDir is called with already existed directory,
	// it should not return error.
	err = maybeCreateWorkingDir(notExistedPath)
	require.NoError(t, err)

	// Prepare dummy file and call maybeCreateWorkingDir.
	// It should return error.
	dummyFilePath := path.Join(tempDir, "dummy_file")
	err = ioutil.WriteFile(dummyFilePath, []byte("dummy_content"), 0644)
	require.NoError(t, err)
	err = maybeCreateWorkingDir(dummyFilePath)
	require.Error(t, err)
}
