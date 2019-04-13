package tests

import (
	"fmt"
	"os"
	"testing"

	"miga/driver"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateCommand(t *testing.T) {

	Convey("All drivers should create migrations files without db config", t, func() {
		dir := "./tmp_crete_test"
		err := os.Mkdir(dir, 0755)
		So(err, ShouldBeNil)
		defer func() {
			err = os.RemoveAll(dir)
			So(err, ShouldBeNil)
		}()

		for driverName, dialects := range drivers {
			for _, dialect := range dialects {
				Convey(fmt.Sprintf("%s driver %s dialect", driverName, dialect), func() {
					cfg := &driver.Config{
						Name:             string(driverName),
						Dialect:          dialect,
						Dir:              dir,
						VersionTableName: string(driverName) + "_db_version",
					}

					driverInst, err := driver.New(cfg)
					So(err, ShouldBeNil)

					err = driverInst.Create(fmt.Sprintf("%s_%s", driverName, dialect), "sql")
					So(err, ShouldBeNil)
				})
			}
		}
	})
}
