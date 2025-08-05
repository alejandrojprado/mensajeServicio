package logger

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var baseLogger log.Logger

func Init() {
	baseLogger = log.NewJSONLogger(os.Stdout)
	baseLogger = log.With(baseLogger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logLevel := os.Getenv("LOG_LEVEL")
	var lvl level.Option
	switch logLevel {
	case "debug":
		lvl = level.AllowDebug()
	case "info":
		lvl = level.AllowInfo()
	case "error":
		lvl = level.AllowError()
	default:
		lvl = level.AllowInfo()
	}
	baseLogger = level.NewFilter(baseLogger, lvl)
}

func LogInfo(msg string, keyvals ...interface{}) {
	_ = level.Info(baseLogger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}

func LogError(msg string, keyvals ...interface{}) {
	_ = level.Error(baseLogger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}

func LogDebug(msg string, keyvals ...interface{}) {
	_ = level.Debug(baseLogger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}
