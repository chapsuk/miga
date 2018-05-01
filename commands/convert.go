package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

func Convert(ctx *cli.Context, cfg *driver.Config) error {
	from := ctx.Args().Get(0)
	if len(from) == 0 {
		return errors.New("FROM_FORAMT required")
	}

	to := ctx.Args().Get(1)
	if len(to) == 0 {
		return errors.New("TO_FORAMT required")
	}

	dest := ctx.Args().Get(2)
	if len(dest) == 0 {
		return errors.New("DESTENITION_PATH required")
	}

	if !driver.Available(from) {
		return fmt.Errorf("unsupported FROM_FORMAT: %s", from)
	}

	if !driver.Available(to) {
		return fmt.Errorf("unsupported TO_FORMAT: %s", from)
	}

	_, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.Mkdir(dest, os.ModeDir)
		if err != nil {
			return err
		}
	}

	return nil
}
