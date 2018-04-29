package driver

import (
	"errors"

	"github.com/chapsuk/miga/driver/goose"
	"github.com/chapsuk/miga/driver/migrate"
)

type Config struct {
	Name string

	Dialect          string
	Dsn              string
	Dir              string
	VersionTableName string
}

type Interface interface {
	Create(name, ext string) error
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
	case "goose":
		return goose.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Dir,
		)
	case "migrate":
		return migrate.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Dir,
		)
	default:
		return nil, errors.New("unsupported driver")
	}
}
