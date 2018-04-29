package config

import (
	"errors"

	"github.com/chapsuk/miga/driver"
	"github.com/spf13/viper"
)

var migrateConfig, seedConfig *driver.Config

func Init(appName, cfg string) error {
	viper.SetConfigFile(cfg)
	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	migrateConfig = &driver.Config{
		Name:             "goose",
		VersionTableName: "miga_db_version",
		Dir:              "./migrations",
	}

	seedConfig = &driver.Config{
		Name:             "goose",
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

	if viper.IsSet("postgres") {
		migrateConfig.Dialect = "postgres"
		migrateConfig.Dsn = viper.GetString("postgres.dsn")
		seedConfig.Dialect = "postgres"
		seedConfig.Dsn = viper.GetString("postgres.dsn")
		return nil
	}

	if viper.IsSet("mysql") {
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
