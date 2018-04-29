package commands

import (
	"errors"

	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

func Create(ctx *cli.Context, d driver.Interface) error {
	name := ctx.Args().Get(0)
	if len(name) == 0 {
		return errors.New("NAME required")
	}

	ext := ctx.Args().Get(1)
	switch ext {
	case "sql":
	case "go":
	default:
		ext = "sql"
	}

	return d.Create(name, ext)
}
