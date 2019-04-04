package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"miga/commands"
	"miga/config"
	"miga/driver"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/urfave/cli.v2"
)

func TestConfig(t *testing.T) {
	Convey("Initialized default config", t, func() {
		err := config.Init("", "", "")
		So(err, ShouldBeNil)

		migCfg := config.MigrateDriverConfig()
		sedCfg := config.SeedDriverConfig()

		Convey("Default driver goose", func() {
			So(migCfg.Name, ShouldEqual, "goose")
			So(sedCfg.Name, ShouldEqual, "goose")
		})

		Convey("Default dialect postgres", func() {
			So(migCfg.Dialect, ShouldEqual, "postgres")
			So(sedCfg.Dialect, ShouldEqual, "postgres")
		})

		Convey("DSN should empty", func() {
			So(migCfg.Dsn, ShouldBeEmpty)
			So(sedCfg.Dsn, ShouldBeEmpty)
		})

		dir := "./migatmp"
		err = os.Mkdir(dir, 0755)
		So(err, ShouldBeNil)

		for driverName := range drivers {
			legend := fmt.Sprintf("Create tmp path setup %s driver", driverName)
			Convey(legend, func() {
				dir += "/" + string(driverName)
				migCfg.Dir = dir + "/migrations"
				sedCfg.Dir = dir + "/seeds"

				for _, d := range []string{dir, migCfg.Dir, sedCfg.Dir} {
					err = os.Mkdir(d, 0755)
					So(err, ShouldBeNil)
				}

				d, err := driver.New(migCfg)
				So(err, ShouldBeNil)

				s, err := driver.New(sedCfg)
				So(err, ShouldBeNil)

				Convey("Miga should complete create cmd with empty DSN", func() {
					flags := &flag.FlagSet{}
					flags.Parse([]string{"testname"})
					ctx := cli.NewContext(nil, flags, nil)

					err = commands.Create(ctx, d)
					So(err, ShouldBeNil)

					err = commands.Create(ctx, s)
					So(err, ShouldBeNil)

					err = os.RemoveAll(dir)
					So(err, ShouldBeNil)
				})

				Convey("Miga should fail up cmd with empty DSN", func() {
					ctx := cli.NewContext(nil, &flag.FlagSet{}, nil)

					err = commands.Up(ctx, d)
					So(err, ShouldNotBeNil)

					err = commands.Up(ctx, s)
					So(err, ShouldNotBeNil)
				})
			})
		}

		err = os.RemoveAll("./migatmp")
		So(err, ShouldBeNil)
	})
}
