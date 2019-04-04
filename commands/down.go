package commands

import (
	"miga/driver"

	"gopkg.in/urfave/cli.v2"
)

// Down rollback last migration
func Down(ctx *cli.Context, d driver.Interface) error {
	return d.Down()
}

// DownTo rollback migrations one by one from current until version from command args
func DownTo(ctx *cli.Context, d driver.Interface) error {
	version, err := parseVersion(ctx)
	if err != nil {
		return err
	}

	return d.DownTo(version)
}
