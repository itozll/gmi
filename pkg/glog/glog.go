package glog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type glog struct {
	l *zap.Logger
	f []zap.Field
}

func New(opt *Options) Logger {
	if opt == nil {
		opt = DefaultOptions()
	}

	l := zap.New(
		zapcore.NewCore(opt.getEncoder(), opt.getWriteSyncer(), opt.getLevel()),
		opt.getAddCaller()...,
	)

	return (&glog{
		l: l,
	})
}

func Attach(l *zap.Logger) Logger {
	return &glog{l: l}
}

func (g *glog) AddCallerSkip(n int) Logger {
	l := g.l.WithOptions(zap.AddCallerSkip(n))
	return &glog{l: l, f: g.f}
}

func (g *glog) WithFields(key string, val interface{}, kvs ...interface{}) Logger {
	f := append(g.f, zap.Any(key, val))

	if len(kvs) > 0 {
		for i := 0; i < len(kvs)-1; i++ {
			key, ok := kvs[i].(string)
			if !ok {
				g.AddCallerSkip(1).WithFields("_item_key", kvs[i], "_err_message", "param must be string").Warn("LogError")
				break
			}

			f = append(f, zap.Any(key, kvs[i+1]))
			i++
		}
	}

	return &glog{l: g.l, f: f}
}

func (g *glog) WithMap(fields map[string]interface{}) Logger {
	if fields == nil {
		return g
	}

	f := g.f[:]

	for key, field := range fields {
		f = append(f, zap.Any(key, field))
	}

	return &glog{l: g.l, f: f}
}

func (g *glog) Debug(msg string) { g.l.Debug(msg, g.f...) }
func (g *glog) Info(msg string)  { g.l.Info(msg, g.f...) }
func (g *glog) Warn(msg string)  { g.l.Warn(msg, g.f...) }
func (g *glog) Error(msg string) { g.l.Error(msg, g.f...) }
func (g *glog) Fatal(msg string) { g.l.Fatal(msg, g.f...) }
func (g *glog) Panic(msg string) { g.l.Panic(msg, g.f...) }
