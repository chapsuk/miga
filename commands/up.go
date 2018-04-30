package commands

import (
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

// Up to latest available migration
func Up(ctx *cli.Context, d driver.Interface) error {
	return d.Up()
}

// UpTo up to version from command args
func UpTo(ctx *cli.Context, d driver.Interface) error {
	version, err := parseVersion(ctx)
	if err != nil {
		return err
	}
	return d.UpTo(version)
}
