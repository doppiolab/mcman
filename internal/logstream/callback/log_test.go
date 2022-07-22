package callback

import (
	"bytes"
	"testing"

	"github.com/doppiolab/mcman/internal/logstream"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogCallback(t *testing.T) {
	out := &bytes.Buffer{}
	logger := zerolog.New(out)
	logBlock := &logstream.LogBlock{
		ChanId: "test",
		Msg:    "dummy text",
	}
	callback := NewLogCallback(logger)

	err := callback(logBlock)

	assert.NoError(t, err)
	assert.Equal(t, out.String(), "{\"level\":\"info\",\"message\":\"[test] dummy text\"}\n")
}
