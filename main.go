package main

import (
	"os"

	"miga/commands/all"
	"miga/commands/migrate"
	"miga/commands/seed"
	"miga/config"
	"miga/logger"

	"gopkg.in/urfave/cli.v2"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/vertica/vertica-sql-go"
)

var (
	Name    = "miga"
	Version = "develop"
)

func main() {
	app := cli.App{
		Name:    Name,
		Version: Version,
		Usage:   "Single CLI for several packages of migration ",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "",
				Usage:   "Config file name",
				EnvVars: []string{"MIGA_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "driver",
				Value:   "",
				Usage:   "Migration driver name: goose, migrate",
				EnvVars: []string{"MIGA_DRIVER"},
			},
			&cli.StringFlag{
				Name:    "log.level",
				Value:   "debug",
				Usage:   "Logger level [debug|info|...]",
				EnvVars: []string{"MIGA_LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log.format",
				Value:   "console",
				Usage:   "Logger output format console|json",
				EnvVars: []string{"MIGA_LOG_FORMAT"},
			},
		},
		Before: func(ctx *cli.Context) error {
			err := logger.Init(
				ctx.App.Name,
				ctx.App.Version,
				ctx.String("log.level"),
				ctx.String("log.format"),
			)
			if err != nil {
				panic("Init logger error: " + err.Error())
			}

			return config.Init(ctx.App.Name, ctx.String("config"), ctx.String("driver"))
		},
		Commands: []*cli.Command{
			all.Command(),
			migrate.Command(),
			seed.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.G().Fatal(err)
	}
}
