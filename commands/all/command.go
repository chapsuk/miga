package all

import (
	"github.com/chapsuk/miga/commands"
	"github.com/chapsuk/miga/config"
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

var (
	migrator driver.Interface
	seeder   driver.Interface
)

// Command returns migration CLI command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "all",
		Usage: "All command combine migration and seed command",
		Before: func(ctx *cli.Context) (err error) {
			migrator, err = driver.New(config.MigrateDriverConfig())
			if err != nil {
				return
			}
			seeder, err = driver.New(config.SeedDriverConfig())
			return
		},
		Subcommands: cli.CommandsByName([]*cli.Command{
			&cli.Command{
				Name:  "up",
				Usage: "Up db to latest migration version and to latest seed.",
				Action: func(ctx *cli.Context) error {
					err := commands.Up(ctx, migrator)
					if err != nil {
						return err
					}
					return commands.Up(ctx, seeder)
				},
			},
		}),
	}
}
