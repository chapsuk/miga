package commands

import (
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

// Create migrations files with given name and extension
func Create(driver func() driver.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create migration file",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 || len(args[0]) == 0 {
				logger.G().Fatalf("File name required")
			}
			name := args[0]

			ext := "sql"
			if len(args) == 2 {
				ext = args[1]
			}

			if err := driver().Create(name, ext); err != nil {
				logger.G().Errorf("create: %s", err)
			}
		},
	}
}
