package goose

import "github.com/chapsuk/miga/logger"

type Goose struct{}

func New(dialect, dsn, tableName, path string) *Goose {
	return &Goose{}
}
func (g *Goose) Create(name, ext string) error {
	logger.G().Infof("name: %s ext: %s", name, ext)
	return nil
}
func (g *Goose) Down() error {
	return nil
}
func (g *Goose) DownTo(version int) error {
	return nil
}
func (g *Goose) Redo() error {
	return nil
}
func (g *Goose) Reset() error {
	return nil
}
func (g *Goose) Status() error {
	return nil
}
func (g *Goose) Up() error {
	return nil
}
func (g *Goose) UpTo(version int) error {
	return nil
}
func (g *Goose) Version() error {
	return nil
}
