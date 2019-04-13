package commands

import (
	"miga/driver"

	"gopkg.in/urfave/cli.v2"
)

// Version print current db version
func Version(ctx *cli.Context, d driver.Interface) error {
	return d.Version()
}
