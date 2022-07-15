package callback

import (
	"github.com/doppiolab/mcman/internal/minecraft/logstream"
	"github.com/rs/zerolog"
)

// LogCallback writes log messages to a zerolog logger.
func NewLogCallback(logger zerolog.Logger) func(*logstream.LogBlock) error {
	return func(l *logstream.LogBlock) error {
		logger.Info().Msg(l.String())
		return nil
	}
}
