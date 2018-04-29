package driver

import (
	"errors"

	"github.com/chapsuk/miga/driver/goose"
)

type Config struct {
	Name string

	Dialect          string
	Dsn              string
	Path             string
	VersionTableName string
}

type Interface interface {
	Create(name, ext string) error
	Down() error
	DownTo(version int) error
	Redo() error
	Reset() error
	Status() error
	Up() error
	UpTo(version int) error
	Version() error
}

func New(cfg *Config) (Interface, error) {
	switch cfg.Name {
	case "goose":
		return goose.New(
			cfg.Dialect,
			cfg.Dsn,
			cfg.VersionTableName,
			cfg.Path,
		), nil
	default:
	}
	return nil, errors.New("unsupported driver")
}
