package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type AppConf struct {
	AppId      string       `mapstructure:"app_id"`
	AppVersion string       `mapstructure:"app_version"`
	Server     ServerConfig `mapstructure:"server"`
	Log        LogConfig    `mapstructure:"log"`
	Worker     WorkerConfig `mapstructure:"worker"`
}

type ServerConfig struct {
	Port int        `mapstructure:"port"`
	CORS CORSConfig `mapstructure:"cors"`
}

type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	MaxAge           int      `mapstructure:"max_age"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json or console
	Folder     string `mapstructure:"folder"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
}

type WorkerConfig struct {
	Count     int `mapstructure:"count"`
	QueueSize int `mapstructure:"queue_size"`
}

var (
	sharedConf *AppConf
	confOnce   sync.Once
)

// SharedAppConf returns the singleton configuration instance
func SharedAppConf() *AppConf {
	confOnce.Do(func() {
		sharedConf = &AppConf{}
	})
	return sharedConf
}

func (c *AppConf) Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}
