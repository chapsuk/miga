package driver

import (
	"errors"

	"github.com/chapsuk/miga/driver/goose"
	"github.com/chapsuk/miga/driver/migrate"
	"github.com/chapsuk/miga/driver/stump"
)

const (
	Goose   = "goose"
	Migrate = "migrate"
	Stump   = "stump"
)

func Available(name string) bool {
	switch name {
	case Goose:
	case Migrate:
	// case Stump:
	default:
		return false
	}
	return true
}

type Config struct {
	Name string

	Dialect          string
	Dsn              string
	Dir              string
	VersionTableName string
}

type Interface interface {
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
	case Stump:
		return stump.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Dir,
		)
	default:
		return nil, errors.New("unsupported migrations driver")
	}
}
