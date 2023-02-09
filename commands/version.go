package commands

import (
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Version print current db version
func Version(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "get current db version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := driver().Version(); err != nil {
				logger.G().Fatalf("get version: %s", err)
			}
		},
	}
}
