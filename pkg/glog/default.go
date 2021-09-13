package glog

var (
	dflOptions = Options{}
	dfl        = New(&dflOptions)
)

func SetDefaultLog(l Logger) { dfl = l }
func GetDefaultLog() Logger  { return dfl }

func WithMap(fields map[string]interface{}) Logger {
	return dfl.WithMap(fields)
}

func WithFields(key string, val interface{}, kvs ...interface{}) Logger {
	return dfl.WithFields(key, val, kvs...)
}

func Debug(msg string) { dfl.AddCallerSkip(1).Debug(msg) }
func Info(msg string)  { dfl.AddCallerSkip(1).Info(msg) }
func Warn(msg string)  { dfl.AddCallerSkip(1).Warn(msg) }
func Error(msg string) { dfl.AddCallerSkip(1).Error(msg) }
func Fatal(msg string) { dfl.AddCallerSkip(1).Fatal(msg) }
func Panic(msg string) { dfl.AddCallerSkip(1).Panic(msg) }
