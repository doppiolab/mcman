package config

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MCMAN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigType("yaml")
}

// Initialize Configuration on startup.
func MustGetConfig(filename string) *Config {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	return config
}
