package commands

import (
	"miga/config"
	"miga/converter"
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"
)

func Convert(cfg func() *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "convert",
		Short: "convert between drivers format",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				logger.G().Fatalf("The up-to version is not defined")
			}
			from := args[0]
			if len(from) == 0 {
				logger.G().Fatal("FROM_FORAMT required")
			}
			to := args[1]
			if len(to) == 0 {
				logger.G().Fatal("TO_FORAMT required")
			}
			dest := args[2]
			if len(dest) == 0 {
				logger.G().Fatal("DESTENITION_PATH required")
			}

			if !driver.Available(from) {
				logger.G().Fatalf("unsupported FROM_FORMAT: %s", from)
			}
			if !driver.Available(to) {
				logger.G().Fatalf("unsupported TO_FORMAT: %s", from)
			}

			if err := converter.Convert(from, to, cfg().Miga.Path, dest); err != nil {
				logger.G().Fatalf("converter: %s", err)
			}
		},
	}
}
