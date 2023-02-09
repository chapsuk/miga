package driver

import (
	"errors"

	"miga/config"
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

func New(cfg *config.Config) (Interface, error) {
	switch cfg.Miga.Driver {
	case Goose:
		return goose.New(
			cfg.Database.Dialect,
			cfg.Database.DSN,
			cfg.Miga.TableName,
			cfg.Miga.Path,
		)
	case Migrate:
		return migrate.New(
			cfg.Database.Dialect,
			cfg.Database.DSN,
			cfg.Miga.TableName,
			cfg.Miga.Path,
		)
	case Impg:
		return impg.New(
			cfg.Database.Dialect,
			cfg.Database.DSN,
			cfg.Miga.TableName,
			cfg.Miga.Path,
		)
	default:
		return nil, errors.New("unsupported migrations driver")
	}
}
