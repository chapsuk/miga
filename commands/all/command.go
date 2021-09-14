package all

import (
	"miga/commands"
	"miga/config"
	"miga/driver"
	"miga/logger"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

// Command returns migration CLI command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "all",
		Usage: "All command combine migration and seed command",
		Subcommands: cli.CommandsByName([]*cli.Command{
			{
				Name:  "up",
				Usage: "Up db to latest migration version and to latest seed.",
				Action: func(ctx *cli.Context) error {
					mcfg := config.MigrateDriverConfig()

					if mcfg.Enabled {
						migrator, err := driver.New(mcfg)
						if err != nil {
							return errors.Wrap(err, "failed create migrator instance")
						}
						err = commands.Up(ctx, migrator)
						if err != nil {
							return errors.Wrap(err, "failed up migrations")
						}
					} else {
						logger.G().Warn("Skip migrate up, migrate dir not exists")
					}

					scfg := config.SeedDriverConfig()
					if scfg.Enabled {
						seeder, err := driver.New(scfg)
						if err != nil {
							return errors.Wrap(err, "failed create seeder instance")
						}
						err = commands.Up(ctx, seeder)
						if err != nil {
							return errors.Wrap(err, "failed up seeds")
						}
					} else {
						logger.G().Warn("Skip seed up, seed dir not exists")
					}

					return nil
				},
			},
		}),
	}
}
