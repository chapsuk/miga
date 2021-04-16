package tests

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"miga/converter"
	"miga/driver"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConvert(t *testing.T) {
	for driverName, dialects := range drivers {
		for _, dialect := range dialects {
			if dialect == "clickhouse" || dialect == "vertica" {
				continue
			}
			for tdriverName := range drivers {
				if tdriverName == driverName {
					continue
				}

				sdir := "./migrations/" + string(tdriverName)
				tdir := "./migrations/tmp/" + string(driverName)

				legend := fmt.Sprintf("Convert from %s to %s source: %s dest: %s",
					strings.ToUpper(string(tdriverName)),
					strings.ToUpper(string(driverName)),
					sdir, tdir)

				Convey(legend, t, func() {
					err := os.MkdirAll(tdir, 0755)
					So(err, ShouldBeNil)

					err = converter.Convert(
						string(tdriverName),
						string(driverName),
						sdir, tdir)
					So(err, ShouldBeNil)

					desc := fmt.Sprintf("Given migrations %s driver with %s dialect",
						strings.ToUpper(string(driverName)),
						strings.ToUpper(dialect),
					)

					Convey(desc, func() {
						driverInst, err := driver.New(&driver.Config{
							Name:             string(driverName),
							Dialect:          dialect,
							Dsn:              string(dsns[dialect]),
							Dir:              tdir,
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

					err = os.RemoveAll(tdir)
					So(err, ShouldBeNil)
				})
			}
		}
	}
}
