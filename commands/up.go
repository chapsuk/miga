package commands

import (
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Up to latest available migration
func Up(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "up db to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := driver().Up(); err != nil {
				logger.G().Errorf("get version: %s", err)
			}
		},
	}
}

// UpTo up to version from command args
func UpTo(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "up-to",
		Short: "up db to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logger.G().Fatalf("The up-to version is not defined")
			}
			version := args[0]
			if err := driver().UpTo(version); err != nil {
				logger.G().Errorf("get version: %s", err)
			}
		},
	}
}
