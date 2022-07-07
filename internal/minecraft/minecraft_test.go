package minecraft

import (
	"fmt"
	"strings"
	"testing"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateMinecraftCommand(t *testing.T) {
	tests := []struct {
		cfg      *config.MinecraftConfig
		expected string
	}{
		{
			cfg: &config.MinecraftConfig{
				JavaCommand: "java",
				JarPath:     "server.jar",
				Args:        []string{"nogui"},
			},
			expected: "java  -jar server.jar nogui",
		},
		{
			cfg: &config.MinecraftConfig{
				JavaCommand: "java13",
				JavaOptions: []string{"-Xms1024m"},
				JarPath:     "server2.jar",
				Args:        []string{"nogui"},
			},
			expected: "java13 -Xms1024m -jar server2.jar nogui",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			actual := createMinecraftCommand(tt.cfg)
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
