package main

import (
	"os"

	"github.com/chapsuk/miga/commands/migrate"
	"github.com/chapsuk/miga/commands/seed"
	"github.com/chapsuk/miga/config"
	"github.com/chapsuk/miga/logger"
	"gopkg.in/urfave/cli.v2"
)

var (
	Name    = "miga"
	Version = "develop"
)

func main() {
	app := cli.App{
		Name:    Name,
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log.level",
				Value:   "debug",
				Usage:   "Logger level [debug|info|...]",
				EnvVars: []string{"MIGA_LOG.LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log.format",
				Value:   "console",
				Usage:   "Logger output format console|json",
				EnvVars: []string{"MIGA_LOG.FORMAT"},
			},
			&cli.StringFlag{
				Name:    "config",
				Value:   "miga.yml",
				Usage:   "Config file name",
				EnvVars: []string{"MIGA_CONFIG"},
			},
		},
		Before: initGlobalsFunc(),
		Commands: []*cli.Command{
			migrate.Command(),
			seed.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.G().Error(err)
	}
}

func initGlobalsFunc() func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		err := config.Init(ctx.App.Name, ctx.String("config"))
		if err != nil {
			return err
		}

		return logger.Init(
			ctx.App.Name,
			ctx.App.Version,
			ctx.String("log.level"),
			ctx.String("log.format"),
		)
	}
}
