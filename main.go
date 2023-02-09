package main

import (
	"log"

	"miga/commands"
	"miga/config"
	"miga/driver"
	"miga/logger"

	"github.com/spf13/cobra"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/vertica/vertica-sql-go"
)

var (
	Name    = "miga"
	Version = "develop"
)

func main() {
	var (
		configFile      string
		migrationDriver driver.Interface
		driverGet       = func() driver.Interface {
			return migrationDriver
		}
		cfg    *config.Config
		cfgGet = func() *config.Config {
			return cfg
		}
		err error
	)

	rootCmd := &cobra.Command{
		Use:   "miga",
		Short: "miga is database migration tool",
		Long:  "miga is database migration tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg, err = config.NewConfig(configFile)
			if err != nil {
				panic(err)
			}
			if err = logger.Init(Name, Version, cfg.Logger.Level, cfg.Logger.Format); err != nil {
				panic(err)
			}

			migrationDriver, err = driver.New(cfg)
			if err != nil {
				logger.G().Fatal("Create driver: %s", err)
			}
			log.Printf("D: %v", migrationDriver)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.G().Info("Root command")
		},
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "miga.yml", "config file path")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(commands.Version(driverGet))
	rootCmd.AddCommand(commands.Up(driverGet))
	rootCmd.AddCommand(commands.UpTo(driverGet))
	rootCmd.AddCommand(commands.Convert(cfgGet))
	rootCmd.AddCommand(commands.Create(driverGet))
	rootCmd.AddCommand(commands.Down(driverGet))
	rootCmd.AddCommand(commands.DownTo(driverGet))
	rootCmd.AddCommand(commands.Redo(driverGet))
	rootCmd.AddCommand(commands.Reset(driverGet))
	rootCmd.AddCommand(commands.Status(driverGet))

	if err := rootCmd.Execute(); err != nil {
		log.Printf("error: %s", err)
	}
}
