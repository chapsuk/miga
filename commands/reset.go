package commands

import (
	"bufio"
	"os"
	"strings"

	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Reset rollback all migrations
func Reset(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "reset database",
		Run: func(cmd *cobra.Command, args []string) {

			logger.G().Info("Rollback all migrations! Are you sure? (yes/no):")
			scanner := bufio.NewScanner(os.Stdin)
			retries := 2
			for scanner.Scan() && retries > 0 {
				ans := strings.ToLower(scanner.Text())
				if ans == "no" {
					logger.G().Info("good choice")
					return
				}
				if ans != "yes" {
					retries--
					logger.G().Info("Enter `yes` or `no`:")
					continue
				}

				break
			}
			logger.G().Info("ok")

			if err := driver().Reset(); err != nil {
				logger.G().Errorf("reset: %s", err)
			}
		},
	}
}
