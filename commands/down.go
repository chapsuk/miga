package commands

import (
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Down rollback last migration
func Down(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "down last migration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := driver().Down(); err != nil {
				logger.G().Errorf("down: %s", err)
			}
		},
	}
}

// DownTo rollback migrations one by one from current until version from command args
func DownTo(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "down-to",
		Short: "down last migration",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logger.G().Fatalf("down-to version is not defined")
			}
			version := args[0]
			if err := driver().DownTo(version); err != nil {
				logger.G().Errorf("down: %s", err)
			}
		},
	}
}
