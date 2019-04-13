package commands

import (
	"miga/driver"

	"gopkg.in/urfave/cli.v2"
)

// Redo rollback and rerun last migration
func Redo(ctx *cli.Context, d driver.Interface) error {
	return d.Redo()
}
