package glog

import "sync"

var (
	dflOptions = Options{}
	dfl        = New(&dflOptions)
	lk         sync.RWMutex
)

func SetDefaultLog(l Logger) { lk.Lock(); defer lk.Unlock(); dfl = l }
func GetDefaultLog() Logger  { lk.RLock(); defer lk.RUnlock(); return dfl }

func WithMap(fields map[string]interface{}) Logger {
	return dfl.WithMap(fields)
}

func WithFields(key string, val interface{}, kvs ...interface{}) Logger {
	return dfl.WithFields(key, val, kvs...)
}

func Debug(msg string) { GetDefaultLog().Debug(msg) }
func Info(msg string)  { GetDefaultLog().Info(msg) }
func Warn(msg string)  { GetDefaultLog().Warn(msg) }
func Error(msg string) { GetDefaultLog().Error(msg) }
func Fatal(msg string) { GetDefaultLog().Fatal(msg) }
func Panic(msg string) { GetDefaultLog().Panic(msg) }
