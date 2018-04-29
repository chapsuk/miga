package commands

import (
	"errors"
	"strconv"

	"gopkg.in/urfave/cli.v2"
)

func parseVersion(ctx *cli.Context) (int, error) {
	versionArg := ctx.Args().Get(0)
	if len(versionArg) == 0 {
		return 0, errors.New("VERSION required")
	}

	version, err := strconv.Atoi(versionArg)
	if err != nil {
		return 0, errors.New("invalid version, should int")
	}
	return version, nil
}
