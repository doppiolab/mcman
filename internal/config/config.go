package config

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Root Config structure.
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Minecraft MinecraftConfig `mapstructure:"minecraft"`
}

type ServerConfig struct {
	// Hostname and Port to listen. Example: 0.0.0.0:8000
	Host string `mapstructure:"host"`
	// Print debug messages if true.
	Debug bool `mapstructure:"debug"`
}

type MinecraftConfig struct {
	JavaCommand string   `mapstructure:"java_command"`
	JarPath     string   `mapstructure:"jar_path"`
	WorkingDir  string   `mapstructure:"working_dir"`
	JavaOptions []string `mapstructure:"java_options"`
	Args        []string `mapstructure:"args"`
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MCMAN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigType("yaml")

	viper.SetDefault("server.host", ":8000")
	viper.SetDefault("server.debug", false)

	viper.SetDefault("minecraft.java", "java")
	viper.SetDefault("minecraft.working_dir", ".")
	viper.SetDefault("minecraft.args", []string{"nogui"})
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
