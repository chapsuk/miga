package config

import (
	"errors"
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
	}

	seedConfig = &driver.Config{
		Name:             driverName,
		VersionTableName: "miga_seed_version",
		Dir:              "./seeds",
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
	if viper.IsSet("seed.table_name") {
		seedConfig.Dir = viper.GetString("seed.path")
	}

	if viper.IsSet("postgres.dsn") || viper.IsSet("postgres.host") {
		dsn := viper.GetString("postgres.dsn")
		if dsn == "" && viper.IsSet("postgres.host") {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
				viper.GetString("postgres.user"),
				viper.GetString("postgres.password"),
				viper.GetString("postgres.host"),
				viper.GetInt("postgres.port"),
				viper.GetString("postgres.db"),
				viper.GetString("postgres.options"),
			)
		}

		migrateConfig.Dialect = "postgres"
		migrateConfig.Dsn = dsn
		seedConfig.Dialect = "postgres"
		seedConfig.Dsn = dsn

		return nil
	}

	if viper.IsSet("mysql.dsn") {
		migrateConfig.Dialect = "mysql"
		migrateConfig.Dsn = viper.GetString("mysql.dsn")
		seedConfig.Dialect = "mysql"
		seedConfig.Dsn = viper.GetString("mysql.dsn")
		return nil
	}

	return errors.New("DB config not found")
}

func MigrateDriverConfig() *driver.Config {
	return migrateConfig
}

func SeedDriverConfig() *driver.Config {
	return seedConfig
}
