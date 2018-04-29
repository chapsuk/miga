package commands

import (
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

func Down(ctx *cli.Context, d driver.Interface) error {
	return d.Down()
}

func DownTo(ctx *cli.Context, d driver.Interface) error {
	version, err := parseVersion(ctx)
	if err != nil {
		return err
	}

	return d.DownTo(version)
}
