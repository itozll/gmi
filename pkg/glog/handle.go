package glog

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Panic(string)

	AddCallerSkip(int) Logger

	WithMap(fields map[string]interface{}) Logger
	WithFields(string, interface{}, ...interface{}) Logger
}
