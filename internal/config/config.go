package config

import (
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from the specified paths or default locations
// If no paths are provided, it looks in "." and "./config"
func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("PARTNER_SERVER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// config.yaml should exist, because environment variables cannot automatically map to the corresponding key
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
