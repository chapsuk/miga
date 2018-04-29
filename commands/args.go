package commands

import (
	"errors"

	"gopkg.in/urfave/cli.v2"
)

func parseVersion(ctx *cli.Context) (string, error) {
	version := ctx.Args().Get(0)
	if len(version) == 0 {
		return "", errors.New("VERSION required")
	}
	return version, nil
}
