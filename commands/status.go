package commands

import (
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

// Status print current migrations state
func Status(ctx *cli.Context, d driver.Interface) error {
	return d.Status()
}
