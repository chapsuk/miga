package seed

import (
	"github.com/chapsuk/miga/commands"
	"github.com/chapsuk/miga/config"
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

var seeder driver.Interface

func Command() *cli.Command {
	return &cli.Command{
		Name:    "seed",
		Aliases: []string{"s"},
		Usage:   "seed command",
		Before: func(ctx *cli.Context) (err error) {
			seeder, err = driver.New(config.SeedDriverConfig())
			return
		},
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:      "create",
				Usage:     "Creates new seed file with next version",
				ArgsUsage: "NAME [sql|go]",
				Action: func(ctx *cli.Context) error {
					return commands.Create(ctx, seeder)
				},
			},
			&cli.Command{
				Name:  "up",
				Usage: "Seed to the most recent version available",
				Action: func(ctx *cli.Context) error {
					return commands.Up(ctx, seeder)
				},
			},
			&cli.Command{
				Name:  "down",
				Usage: "Roll back last seeds",
				Action: func(ctx *cli.Context) error {
					return commands.Down(ctx, seeder)
				},
			},
			&cli.Command{
				Name:  "reset",
				Usage: "Reset all seeds",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "skip warning dialog",
						Value: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					return commands.Reset(ctx, seeder)
				},
			},
			&cli.Command{
				Name:  "version",
				Usage: "Print the current version of seeds",
				Action: func(ctx *cli.Context) error {
					return commands.Version(ctx, seeder)
				},
			},
		},
	}
}
