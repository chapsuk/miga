package commands

import (
	"bufio"
	"os"
	"strings"

	"github.com/chapsuk/miga/driver"
	"github.com/chapsuk/miga/logger"
	"gopkg.in/urfave/cli.v2"
)

func Reset(ctx *cli.Context, d driver.Interface) error {
	if !ctx.Bool("force") {
		logger.G().Info("Rollback all migrations! Are you sure? (yes/no):")
		scanner := bufio.NewScanner(os.Stdin)
		retries := 2
		for scanner.Scan() && retries > 0 {
			ans := strings.ToLower(scanner.Text())
			if ans == "no" {
				logger.G().Info("good choice")
				return nil
			}
			if ans != "yes" {
				retries--
				logger.G().Info("Enter `yes` or `no`:")
				continue
			}

			break
		}
		logger.G().Info("ok")
	}

	return d.Reset()
}
