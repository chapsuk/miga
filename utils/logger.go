package utils

import "github.com/chapsuk/miga/logger"

type StdLogger struct{}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	logger.G().Infof(format, v...)
}

func (l *StdLogger) Verbose() bool {
	return true
}

func (l *StdLogger) Fatal(v ...interface{}) {
	logger.G().Fatal(v)
}
func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	logger.G().Fatalf(format, v...)
}
func (l *StdLogger) Print(v ...interface{}) {
	logger.G().Info(v...)
}
func (l *StdLogger) Println(v ...interface{}) {
	logger.G().Info(v...)
}
