package tests

import (
	"database/sql"

	"miga/driver"
)

type (
	testCase struct {
		Description string
		Action      func(driver.Interface)
		Assert      func(db *sql.DB, tableSuffix string)
		Condition   func(driverName, dialect string) bool
	}

	dsn        string
	driverName string
	dialects   []string
)

var (
	drivers = map[driverName]dialects{
		"goose":   []string{"mysql", "postgres", "clickhouse", "vertica", "clickhouse-replicated"},
		"migrate": []string{"mysql", "postgres"},
		"impg":    []string{"postgres"},
	}

	dsns = map[string]dsn{
		"postgres":              "postgres://user:password@127.0.0.1:5432/miga?sslmode=disable",
		"mysql":                 "user:password@tcp(127.0.0.1:3306)/miga",
		"clickhouse":            "tcp://user:password@127.0.0.1:9000/miga",
		"clickhouse-replicated": "tcp://user:password@127.0.0.1:9000/miga",
		"vertica":               "vertica://dbadmin:@127.0.0.1:5433/docker",
	}
)
