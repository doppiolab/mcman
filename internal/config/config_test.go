package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustGetConfig(t *testing.T) {
	cfg := MustGetConfig("../../default.yaml")

	assert.Equal(t, cfg.Server.Host, "0.0.0.0:8000")
}
