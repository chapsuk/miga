package seed

import (
	"errors"

	"miga/commands"
	"miga/config"
	"miga/driver"

	"gopkg.in/urfave/cli.v2"
)

var seeder driver.Interface

// Command returns seed CLI command
func Command() *cli.Command {
	return &cli.Command{
		Name:  "seed",
		Usage: "Seeding root command, see",
		Before: func(ctx *cli.Context) (err error) {
			scfg := config.SeedDriverConfig()
			if !scfg.Enabled {
				err = errors.New("seeds dir not exists or misconfigured")
				return
			}
			seeder, err = driver.New(scfg)
			return
		},
		Subcommands: cli.CommandsByName([]*cli.Command{
			&cli.Command{
				Name:      "convert",
				Usage:     "Converting seeds to another format",
				ArgsUsage: "FROM TO DESTENITION_PATH",
				Action: func(ctx *cli.Context) error {
					return commands.Convert(ctx, config.MigrateDriverConfig())
				},
			},
			&cli.Command{
				Name:      "create",
				Usage:     "Creates new seed sql|go file with next version",
				ArgsUsage: "NAME EXT",
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
				Name:      "up-to",
				ArgsUsage: "VERSION",
				Usage:     "Goto a specific seed VERSION",
				Action: func(ctx *cli.Context) error {
					return commands.UpTo(ctx, seeder)
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
				Name:      "down-to",
				ArgsUsage: "VERSION",
				Usage:     "Roll back to a specific seed VERSION",
				Action: func(ctx *cli.Context) error {
					return commands.DownTo(ctx, seeder)
				},
			},
			&cli.Command{
				Name:  "redo",
				Usage: "Re-run the latest seed",
				Action: func(ctx *cli.Context) error {
					return commands.Redo(ctx, seeder)
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
		}),
	}
}
