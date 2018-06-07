package all

import (
	"github.com/chapsuk/miga/commands"
	"github.com/chapsuk/miga/config"
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

// Command returns migration CLI command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "all",
		Usage: "All command combine migration and seed command",
		Subcommands: cli.CommandsByName([]*cli.Command{
			&cli.Command{
				Name:  "up",
				Usage: "Up db to latest migration version and to latest seed.",
				Action: func(ctx *cli.Context) error {
					migrator, err := driver.New(config.MigrateDriverConfig())
					if err != nil {
						return err
					}
					err = commands.Up(ctx, migrator)
					if err != nil {
						return err
					}

					seeder, err := driver.New(config.SeedDriverConfig())
					if err != nil {
						return err
					}
					return commands.Up(ctx, seeder)
				},
			},
		}),
	}
}
