package tests

import (
	"miga/logger"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoggerInit(t *testing.T) {
	Convey("Should init logger with console and json formats without errors", t, func() {
		err := logger.Init("miga", "v0.0.1", "info", "console")
		So(err, ShouldBeNil)

		err = logger.Init("miga", "v0.0.1", "info", "json")
		So(err, ShouldBeNil)
	})
}
