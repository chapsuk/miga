package driver

import (
	"errors"

	"miga/driver/goose"
	"miga/driver/impg"
	"miga/driver/migrate"
)

const (
	Goose   = "goose"
	Migrate = "migrate"
	Impg    = "impg"
)

func Available(name string) bool {
	switch name {
	case Goose:
	case Migrate:
	case Impg:
	default:
		return false
	}
	return true
}

type (
	Config struct {
		Name             string
		Dialect          string
		Dsn              string
		Dir              string
		VersionTableName string
		Enabled          bool

		ClickhouseSchema      string
		ClickhouseClusterName string
		ClickhouseEngine      string
		ClickhouseSharded     bool
	}

	Interface interface {
		Create(name, ext string) error
		Close() error
		Down() error
		DownTo(version string) error
		Redo() error
		Reset() error
		Status() error
		Up() error
		UpTo(version string) error
		Version() error
	}
)

func New(cfg *Config) (Interface, error) {
	switch cfg.Name {
	case Goose:
		return goose.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Dir,
		)
	case Migrate:
		return migrate.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Dir,
		)
	case Impg:
		return impg.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Dir,
		)
	default:
		return nil, errors.New("unsupported migrations driver")
	}
}

func (c *Config) HasDBConfig() bool {
	return c.Dsn != "" && c.Dialect != ""
}
