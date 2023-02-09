package tests

import (
	"testing"

	"miga/config"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Initialized config", t, func() {
		cfg, err := config.NewConfig("./../miga.yml")
		So(err, ShouldBeNil)

		Convey("Should parse miga config", func() {
			So(cfg.Miga.Driver, ShouldEqual, "goose")
			So(cfg.Miga.Path, ShouldEqual, "./tests/migrations/goose")
			So(cfg.Miga.TableName, ShouldEqual, "db_version")
		})

		Convey("Should parse logger", func() {
			So(cfg.Logger.Level, ShouldEqual, "info")
			So(cfg.Logger.Format, ShouldEqual, "console")
		})

		Convey("Should parse db block", func() {
			So(cfg.Database.DSN, ShouldNotBeEmpty)
			So(cfg.Database.Dialect, ShouldEqual, "mysql")
		})
	})
}
