package tests

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/chapsuk/miga/driver"
	. "github.com/smartystreets/goconvey/convey"
)

var migrationCases = []testCase{
	{
		Description: "#1 up to first migartions (create users table)",
		Action: func(d driver.Interface) {
			err := d.UpTo("1")
			So(err, ShouldBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users")
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
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
		Description: "#101 incorrect migration (duplicate of 3 migration)",
		Action: func(d driver.Interface) {
			err := d.UpTo("101")
			So(err, ShouldNotBeNil)
		},
		Assert: func(db *sql.DB) {
			r, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
			_, err := db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
			_, err = db.Query("SELECT COUNT(*) FROM users")
			So(err, ShouldNotBeNil)
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
			_, err = db.Query("SELECT COUNT(*) FROM users")
			So(err, ShouldBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
			_, err := db.Query("SELECT COUNT(*) FROM users")
			So(err, ShouldBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
			_, err = db.Query("SELECT COUNT(*) FROM users")
			So(err, ShouldNotBeNil)
			_, err = db.Query("SELECT COUNT(*) FROM users WHERE email='test'")
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
				driverInst, err := driver.New(&driver.Config{
					Name:             string(driverName),
					Dialect:          dialect,
					Dsn:              string(dsns[dialect]),
					Dir:              "./migrations/" + string(driverName),
					VersionTableName: string(driverName) + "_db_version",
				})
				So(err, ShouldBeNil)

				db, err := sql.Open(dialect, string(dsns[dialect]))
				So(err, ShouldBeNil)

				for _, testCase := range migrationCases {
					Convey(testCase.Description, func() {
						testCase.Action(driverInst)
						testCase.Assert(db)
					})
				}
			})
		}
	}
}
