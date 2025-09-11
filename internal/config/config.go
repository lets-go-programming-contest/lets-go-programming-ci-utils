package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Mode string

const (
	SkipMode    Mode = "skip"
	DefaultMode Mode = "default"
	CommonMode  Mode = "common"
	StudentMode Mode = "student"
)

type StageConfig struct {
	Mode Mode `mapstructure:"mode"`
}

type Config struct {
	BuildMode StageConfig `mapstructure:"build"`
	LintMode  StageConfig `mapstructure:"lint"`
	TestMode  StageConfig `mapstructure:"test"`
}

func ReadConfig(configFile string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read task configuration: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal task configuration: %w", err)
	}

	return cfg, nil
}
