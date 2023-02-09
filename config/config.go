package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Miga     MigaConfig     `yaml:"miga"`
		Logger   LoggerConfig   `yaml:"logger"`
		Database DatabaseConfig `yaml:"db"`
	}

	MigaConfig struct {
		Driver    string `yaml:"driver"`
		Path      string `yaml:"path"`
		TableName string `yaml:"table"`
	}
	DatabaseConfig struct {
		DSN     string `yaml:"dsn"`
		Dialect string `yaml:"dialect"`
	}

	LoggerConfig struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	}
)

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("miga")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{
		Miga: MigaConfig{
			Driver:    viper.GetString("miga.driver"),
			Path:      viper.GetString("miga.path"),
			TableName: viper.GetString("miga.table"),
		},
		Logger: LoggerConfig{
			Level:  viper.GetString("logger.level"),
			Format: viper.GetString("logger.format"),
		},
		Database: DatabaseConfig{
			DSN:     viper.GetString("db.dsn"),
			Dialect: viper.GetString("db.dialect"),
		},
	}

	if len(cfg.Logger.Level) == 0 {
		cfg.Logger.Level = "info"
	}
	if len(cfg.Logger.Format) == 0 {
		cfg.Logger.Format = "console"
	}

	return cfg, nil
}
