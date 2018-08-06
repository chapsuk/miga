package config

import (
	"fmt"
	"strings"

	"github.com/chapsuk/miga/driver"
	"github.com/chapsuk/miga/logger"
	"github.com/spf13/viper"
)

var migrateConfig, seedConfig *driver.Config

// Init configuration with viper
func Init(appName, cfg, driverName string) error {
	viper.SetConfigFile(cfg)
	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		if !strings.Contains(err.Error(), "Not Found") {
			return err
		}
		logger.G().Debugf("Missing config file: %s", err)
	}

	if driverName == "" {
		driverName = viper.GetString("driver")
		if !viper.IsSet("driver") {
			driverName = driver.Goose
		}
	}
	if !driver.Available(driverName) {
		return fmt.Errorf("Unsupported driver %s", driverName)
	}
	logger.G().Infof("Using %s driver", driverName)

	migrateConfig = &driver.Config{
		Name:             driverName,
		VersionTableName: "miga_db_version",
		Dir:              "./migrations",
		Dialect:          "postgres",
	}

	seedConfig = &driver.Config{
		Name:             driverName,
		VersionTableName: "miga_seed_version",
		Dir:              "./seeds",
		Dialect:          "postgres",
	}

	if viper.IsSet("migrate.table_name") {
		migrateConfig.VersionTableName = viper.GetString("migrate.table_name")
	}
	if viper.IsSet("seed.table_name") {
		seedConfig.VersionTableName = viper.GetString("seed.table_name")
	}

	if viper.IsSet("migrate.path") {
		migrateConfig.Dir = viper.GetString("migrate.path")
	}
	if viper.IsSet("seed.path") {
		seedConfig.Dir = viper.GetString("seed.path")
	}

	return nil
}

func MigrateDriverConfig() *driver.Config {
	fillDBConfig(migrateConfig)
	return migrateConfig
}

func SeedDriverConfig() *driver.Config {
	fillDBConfig(seedConfig)
	return seedConfig
}

func fillDBConfig(cfg *driver.Config) {
	if cfg.HasDBConfig() {
		return
	}

	if viper.IsSet("postgres.dsn") || viper.IsSet("postgres.host") {
		cfg.Dialect = "postgres"
		dsn := viper.GetString("postgres.dsn")
		if dsn == "" {
			var addr string
			if viper.IsSet("postgres.host") {
				port := viper.GetInt("postgres.port")
				if port == 0 {
					port = 5432
				}
				addr = fmt.Sprintf("%s:%d", viper.GetString("postgres.host"), port)
			}

			if viper.GetString("postgres.address") != "" {
				addr = viper.GetString("postgres.address")
			}

			db := viper.GetString("postgres.db")
			if viper.IsSet("postgres.database") {
				db = viper.GetString("postgres.database")
			}

			dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?%s",
				viper.GetString("postgres.user"),
				viper.GetString("postgres.password"),
				addr,
				db,
				viper.GetString("postgres.options"),
			)
		}
		cfg.Dsn = dsn
		return
	}

	if viper.IsSet("mysql.dsn") {
		cfg.Dialect = "mysql"
		cfg.Dsn = viper.GetString("mysql.dsn")
		return
	}
}
