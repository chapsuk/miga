package config

import (
	"database/sql"
	"strings"
	"testing"

	"miga/driver"

	"github.com/go-pg/pg"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

type (
	credentials struct {
		Host     string
		Port     string
		Address  string
		User     string
		Password string
		Database string
		Options  string
	}

	testCase struct {
		creds      credentials
		shouldErr  bool
		invalidDSN bool
		dsn        string
	}
)

func TestPostgresDSN(t *testing.T) {
	for i, tcase := range pgCredentialsTable {
		viper.Set("postgres.host", tcase.creds.Host)
		viper.Set("postgres.port", tcase.creds.Port)
		viper.Set("postgres.address", tcase.creds.Address)
		viper.Set("postgres.user", tcase.creds.User)
		viper.Set("postgres.database", tcase.creds.Database)
		viper.Set("postgres.password", tcase.creds.Password)
		viper.Set("postgres.options", tcase.creds.Options)

		dcfg := &driver.Config{}
		fillDBConfig(dcfg)

		if dcfg.Dsn != tcase.dsn {
			t.Fatalf("#%d expected dsn: %s actual: %s", i, tcase.dsn, dcfg.Dsn)
		}

		opts, errPg := pg.ParseURL(dcfg.Dsn)
		if errPg == nil {
			db := pg.Connect(opts)
			_, err := db.Exec("select 1")
			if err != nil {
				if !tcase.invalidDSN && !strings.Contains(err.Error(), "connection refused") {
					t.Fatalf("Unexpected connection error: %s", err)
				}
			}
		}

		uri, errPq := pq.ParseURL(dcfg.Dsn)
		if errPq == nil {
			db, err := sql.Open("postgres", uri)
			if err != nil {
				t.Logf("ERROR %s", err)
			}

			_, err = db.Query("select 1")
			if err != nil {
				if !tcase.invalidDSN && !strings.Contains(err.Error(), "connection refused") {
					t.Fatalf("Unexpected connection error: %s", err)
				}
			}
		}

		errs := 0
		if tcase.shouldErr {
			if errPg == nil {
				errs++
				t.Logf("#%d expected `pg` dsn parse error, but err is nil", i)
			}
			if errPq == nil {
				errs++
				t.Logf("#%d expected `pq` dsn parse error, but err is nil", i)
			}
			if errs > 0 {
				t.Fatalf("expected dsn parse error, but pg: %s pq: %s", errPg, errPq)
			}
		} else {
			if errPg != nil {
				errs++
				t.Logf("#%d `pg` dsn parse error: %s", i, errPg)
			}
			if errPq != nil {
				errs++
				t.Logf("#%d `pq` dsn parse error: %s", i, errPq)
			}
			if errs > 0 {
				t.Fatalf("expected no dsn parse error, but pg %s pq: %s", errPg, errPq)
			}
		}
	}
}

var pgCredentialsTable = []testCase{
	{
		creds: credentials{
			Host:     "127.0.0.1",
			Port:     "11132",
			Address:  "",
			User:     "foo",
			Password: "bar",
			Database: "test",
			Options:  "",
		},
		dsn: "postgres://foo:bar@127.0.0.1:11132/test?",
	},
	{
		creds: credentials{
			Host:     "128.0.0.1",
			Port:     "1331",
			Address:  "127.0.0.1:11132",
			User:     "foo",
			Password: "bar",
			Database: "test",
			Options:  "",
		},
		dsn: "postgres://foo:bar@127.0.0.1:11132/test?",
	},
	{
		creds: credentials{
			Address:  "127.0.0.1:11132",
			User:     "foo",
			Password: "bar/",
			Database: "test",
			Options:  "",
		},
		invalidDSN: true,
		shouldErr:  true,
		dsn:        "postgres://foo:bar/@127.0.0.1:11132/test?",
	},
}
