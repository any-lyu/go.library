package logs

import (
	"os"

	"github.com/astaxie/beego/logs"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

func init() {
	// logs.Async()
	logs.SetLogFuncCall(true)
	logs.SetLevel(logs.LevelDebug)
	logs.SetLogFuncCallDepth(4)
	_ = logs.SetLogger(logs.AdapterConsole)
	//	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"../logs/library_back.logs"}`)
}

// SetLevel .
func SetLevel(level int) {
	logs.SetLevel(level)
}

// Async .
func Async(msgLen ...int64) {
	logs.Async(msgLen...)
}

// SetLogFuncCall .
func SetLogFuncCall(b bool) {
	logs.SetLogFuncCall(b)
}

// SetLogFuncCallDepth .
func SetLogFuncCallDepth(d int) {
	logs.SetLogFuncCallDepth(d)
}

// SetLogger .
func SetLogger(adapter string, config ...string) error {
	return logs.SetLogger(adapter, config...)
}

type newLoggerFunc func() logs.Logger

// Logger .
type Logger logs.Logger

// Register .
func Register(name string, log newLoggerFunc) {
	logs.Register(name, func() logs.Logger {
		return log()
	})
}

//Info log info
func Info(f interface{}, v ...interface{}) {
	logs.Info(f, v...)
}

//Debug log debug
func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
}

//Error log error
func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v...)
}

//Warn log warn
func Warn(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)
}

//Fatal log fatal
func Fatal(v ...interface{}) {
	Error(v)
	os.Exit(1)
}

// GetLevel log level
func GetLevel() int {
	return logs.GetBeeLogger().GetLevel()
}
