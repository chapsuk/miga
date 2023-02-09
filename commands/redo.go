package commands

import (
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Redo rollback and rerun last migration
func Redo(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "redo",
		Short: "redo cmd",
		Run: func(cmd *cobra.Command, args []string) {
			if err := driver().Redo(); err != nil {
				logger.G().Fatalf("redo: %s", err)
			}
		},
	}
}
