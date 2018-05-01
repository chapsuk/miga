package migrate

import (
	"github.com/chapsuk/miga/commands"
	"github.com/chapsuk/miga/config"
	"github.com/chapsuk/miga/driver"
	"gopkg.in/urfave/cli.v2"
)

var migrator driver.Interface

// Command returns migration CLI command
func Command() *cli.Command {
	return &cli.Command{
		Name:    "migrate",
		Aliases: []string{"m"},
		Usage:   "Migrations root command",
		Before: func(ctx *cli.Context) (err error) {
			migrator, err = driver.New(config.MigrateDriverConfig())
			return
		},
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:      "convert",
				Usage:     "Converting migrations FROM_FORMAT to TO_FORMAT and store to DESTENITION_PATH",
				ArgsUsage: "FROM_FORMAT TO_FORMAT DESTENITION_PATH",
				Action: func(ctx *cli.Context) error {
					return commands.Convert(ctx, config.MigrateDriverConfig())
				},
			},
			&cli.Command{
				Name:      "create",
				Usage:     "Creates new migration sql file",
				ArgsUsage: "NAME",
				Action: func(ctx *cli.Context) error {
					return commands.Create(ctx, migrator)
				},
			},
			&cli.Command{
				Name:  "down",
				Usage: "Roll back the version by 1",
				Action: func(ctx *cli.Context) error {
					return commands.Down(ctx, migrator)
				},
			},
			&cli.Command{
				Name:      "down-to",
				ArgsUsage: "VERSION",
				Usage:     "Roll back to a specific VERSION",
				Action: func(ctx *cli.Context) error {
					return commands.DownTo(ctx, migrator)
				},
			},
			&cli.Command{
				Name:  "redo",
				Usage: "Re-run the latest migration",
				Action: func(ctx *cli.Context) error {
					return commands.Redo(ctx, migrator)
				},
			},
			&cli.Command{
				Name:  "reset",
				Usage: "Roll back all migrations",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "skip warning dialog",
						Value: false,
					},
				},
				Action: func(ctx *cli.Context) error {
					return commands.Reset(ctx, migrator)
				},
			},
			&cli.Command{
				Name:  "status",
				Usage: "Dump the migration status for the current DB",
				Action: func(ctx *cli.Context) error {
					return commands.Status(ctx, migrator)
				},
			},
			&cli.Command{
				Name:  "up",
				Usage: "Migrate the DB to the most recent version available",
				Action: func(ctx *cli.Context) error {
					return commands.Up(ctx, migrator)
				},
			},
			&cli.Command{
				Name:      "up-to",
				ArgsUsage: "VERSION",
				Usage:     "Migrate the DB to a specific VERSION",
				Action: func(ctx *cli.Context) error {
					return commands.UpTo(ctx, migrator)
				},
			},
			&cli.Command{
				Name:  "version",
				Usage: "Print the current version of the database",
				Action: func(ctx *cli.Context) error {
					return commands.Version(ctx, migrator)
				},
			},
		},
	}
}
