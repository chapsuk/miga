package goose

import (
	"database/sql"
	"fmt"
)

// OpenDBWithDriver creates a connection to a database, and modifies goose
// internals to be compatible with the supplied driver by calling SetDialect.
func OpenDBWithDriver(driver string, dbstring string) (*sql.DB, error) {
	if err := SetDialect(driver); err != nil {
		return nil, err
	}

	switch driver {
	case "mssql":
		driver = "sqlserver"
	case "redshift":
		driver = "postgres"
	case "tidb":
		driver = "mysql"
	case "clickhouse-replicated":
		db, err := sql.Open("clickhouse", dbstring)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}
		_, err = db.Exec("SET insert_quorum=2")
		if err != nil {
			return nil, fmt.Errorf("set insert_quorum %w", err)
		}
		_, err = db.Exec("SET select_sequential_consistency=1")
		if err != nil {
			return nil, fmt.Errorf("set select_sequential_consistency %w", err)
		}
		return db, nil
	}

	switch driver {
	case "postgres", "pgx", "sqlite3", "sqlite", "mysql", "sqlserver", "clickhouse", "vertica":
		return sql.Open(driver, dbstring)
	default:
		return nil, fmt.Errorf("unsupported driver %s", driver)
	}
}
