package commands

import (
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Status print current migrations state
func Status(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "returns current db status",
		Run: func(cmd *cobra.Command, args []string) {
			if err := driver().Status(); err != nil {
				logger.G().Fatalf("status: %s", err)
			}
		},
	}
}
