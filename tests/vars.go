package tests

import (
	"database/sql"

	"github.com/chapsuk/miga/driver"
)

type (
	testCase struct {
		Description string
		Action      func(driver.Interface)
		Assert      func(db *sql.DB)
	}

	dsn        string
	driverName string
	dialects   []string
)

var (
	drivers = map[driverName]dialects{
		"goose":   []string{"mysql", "postgres"},
		"migrate": []string{"mysql", "postgres"},
	}

	dsns = map[string]dsn{
		"postgres": "host=127.0.0.1 user=user password=password port=5432 sslmode=disable database=miga",
		"mysql":    "user:password@tcp(127.0.0.1:3306)/miga",
	}
)
