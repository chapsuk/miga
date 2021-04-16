package tests

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"miga/driver"

	_ "github.com/ClickHouse/clickhouse-go"
	. "github.com/smartystreets/goconvey/convey"
	_ "github.com/vertica/vertica-sql-go"
)

var migrationCases = []testCase{
	{
		Description: "#1 up to first migartions (create users table)",
		Action: func(d driver.Interface) {
			err := d.UpTo("1")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users where migastas >= 0")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(count, ShouldEqual, 0)
			So(r.Err(), ShouldBeNil)
		},
	},
	{
		Description: "#2 up to second (create wallets table)",
		Action: func(d driver.Interface) {
			err := d.UpTo("2")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM wallets")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(count, ShouldEqual, 0)
			So(r.Err(), ShouldBeNil)
		},
	},
	{
		Description: "#3 up to third (alter tables)",
		Action: func(d driver.Interface) {
			err := d.UpTo("3")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(count, ShouldEqual, 0)
			So(r.Err(), ShouldBeNil)
		},
	},
	{
		Description: "#4 retry last migration",
		Action: func(d driver.Interface) {
			err := d.Redo()
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(count, ShouldEqual, 0)
			So(r.Err(), ShouldBeNil)
		},
	},
	{
		Description: "#5 query with `;`",
		Action: func(d driver.Interface) {
			err := d.UpTo("4")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE name='Abib;Rabib' and migastas >= 0")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(count, ShouldEqual, 1)
			So(r.Err(), ShouldBeNil)
		},
		Condition: func(driverName, dialect string) bool {
			return dialect != "clickhouse"
		},
	},
	{
		Description: "#6 plpsql statement, should create histories table and func for create inheritans",
		Action: func(d driver.Interface) {
			err := d.UpTo("5")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM histories")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(r.Err(), ShouldBeNil)
			So(count, ShouldEqual, 0)

			_, err = db.Query("SELECT histories_partition_creation('now', 'now');")
			So(r.Err(), ShouldBeNil)

			r, err = db.Query(fmt.Sprintf("SELECT COUNT(*) FROM histories_%d_%02d", time.Now().Year(), time.Now().Month()))
			So(err, ShouldBeNil)
			for r.Next() {
				r.Scan(&count)
			}
			So(r.Err(), ShouldBeNil)
			So(count, ShouldEqual, 0)
		},
		Condition: func(driverName, dialect string) bool {
			return dialect == "postgres"
		},
	},
	{
		Description: "#7 (goose_issue158): create custom type",
		Action: func(d driver.Interface) {
			err := d.UpTo("6")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT 1 FROM pg_type WHERE typname = 'things'")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(r.Err(), ShouldBeNil)
			So(count, ShouldEqual, 1)

			_, err = db.Exec("INSERT INTO doge (id, th) VALUES (1, 'hello')")
			So(r.Err(), ShouldBeNil)

			r, err = db.Query("SELECT id, th FROM doge")
			So(err, ShouldBeNil)
			for r.Next() {
				var (
					id int
					th string
				)
				err = r.Scan(&id, &th)
				So(err, ShouldBeNil)
				So(id, ShouldEqual, 1)
				So(th, ShouldEqual, "hello")
			}
			So(r.Err(), ShouldBeNil)
		},
		Condition: func(driverName, dialect string) bool {
			return dialect == "postgres"
		},
	},
	{
		Description: "#101 incorrect migration (duplicate of 3 migration)",
		Action: func(d driver.Interface) {
			err := d.UpTo("101")
			So(err, ShouldNotBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldBeNil)
			count := 1
			for r.Next() {
				r.Scan(&count)
			}
			So(count, ShouldEqual, 0)
			So(r.Err(), ShouldBeNil)

			_, err = db.Query("SELECT COUNT(*) FROM foo")
			So(err, ShouldNotBeNil)
		},
	},
	{
		Description: "#102 try jump over failed migration",
		Action: func(d driver.Interface) {
			err := d.UpTo("102")
			So(err, ShouldNotBeNil)
		},
		Assert: func(db *sql.DB) {
			_, err := db.Query("SELECT COUNT(*) FROM never")
			So(err, ShouldNotBeNil)
		},
	},
	{
		Description: "#103 down to 2 migartion",
		Action: func(d driver.Interface) {
			err := d.DownTo("2")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			_, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldNotBeNil)
		},
	},
	{
		Description: "#104 reset all",
		Action: func(d driver.Interface) {
			err := d.Reset()
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			_, err := db.Query("SELECT COUNT(*) FROM wallets")
			So(err, ShouldNotBeNil)
			// _, err = db.Query("SELECT COUNT(*) FROM users")
			// So(err, ShouldNotBeNil)
		},
	},
	{
		Description: "#105 up to latest, but stop on failed",
		Action: func(d driver.Interface) {
			err := d.Up()
			So(err, ShouldNotBeNil)
		},
		Assert: func(db *sql.DB) {
			_, err := db.Query("SELECT COUNT(*) FROM wallets")
			So(err, ShouldBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users where migastas >= 0")
			So(err, ShouldBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM never")
			So(err, ShouldNotBeNil)
		},
	},
	{
		Description: "#106 down to first",
		Action: func(d driver.Interface) {
			err := d.DownTo("1")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			_, err := db.Query("SELECT COUNT(*) FROM users where migastas >= 0")
			So(err, ShouldBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldNotBeNil)
		},
	},
	{
		Description: "#107 reset all",
		Action: func(d driver.Interface) {
			err := d.Reset()
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			_, err := db.Query("SELECT COUNT(*) FROM wallets")
			So(err, ShouldNotBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users where migastas >= 0")
			So(err, ShouldNotBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users WHERE email='test' and migastas >= 0")
			So(err, ShouldNotBeNil)
		},
	},
}

func TestMigrations(t *testing.T) {
	// logger.Init("miga_test", "test", "info", "console")
	for driverName, dialects := range drivers {
		for _, dialect := range dialects {
			desc := fmt.Sprintf("Given migrations %s driver with %s dialect",
				strings.ToUpper(string(driverName)),
				strings.ToUpper(dialect),
			)
			Convey(desc, t, func() {
				dir := "./migrations/" + string(driverName)
				if dialect == "clickhouse" {
					dir += "_" + dialect
				}

				driverInst, err := driver.New(&driver.Config{
					Name:             string(driverName),
					Dialect:          dialect,
					Dsn:              string(dsns[dialect]),
					Dir:              dir,
					VersionTableName: string(driverName) + "_db_version",
				})
				So(err, ShouldBeNil)
				defer driverInst.Close()

				db, err := sql.Open(dialect, string(dsns[dialect]))
				So(err, ShouldBeNil)
				defer db.Close()

				for _, testCase := range migrationCases {
					if testCase.Condition != nil {
						if !testCase.Condition(string(driverName), dialect) {
							continue
						}
					}

					Convey(testCase.Description, func() {
						testCase.Action(driverInst)
						testCase.Assert(db)
					})
				}
			})
		}
	}
}
