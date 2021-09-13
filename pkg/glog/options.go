package glog

import (
	"os"
	"strconv"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	OutputStdout = "stdout" // default
	OutputStderr = "stderr"
	OutputFile   = "file"

	FormatJSON = "json" // default
	FormatText = "text"

	LevelDebug = "debug" // default
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
)

type Options struct {
	Output     string `json:"output,omitempty" yaml:"output" toml:"output"`
	Format     string `json:"format,omitempty" yaml:"format" toml:"format"`
	Level      string `json:"level,omitempty" yaml:"level" toml:"level"`
	TimeFormat string `json:"time_format,omitempty" yaml:"time_format" toml:"time_format"`

	// Output == OutputFile
	FileName string `json:"file_name,omitempty" yaml:"file_name" toml:"file_name"`
	Size     string `json:"size,omitempty" yaml:"size" toml:"size"`

	NoAddCaller bool `json:"no_caller,omitempty" yaml:"no_caller" toml:"no_caller"`
}

// DefaultOptions default Options
func DefaultOptions() *Options                               { return &Options{} }
func (o *Options) WithOutput(Output string) *Options         { o.Output = Output; return o }
func (o *Options) WithFormat(Format string) *Options         { o.Format = Format; return o }
func (o *Options) WithLevel(Level string) *Options           { o.Level = Level; return o }
func (o *Options) WithTimeFormat(TimeFormat string) *Options { o.TimeFormat = TimeFormat; return o }
func (o *Options) WithFileName(FileName string) *Options     { o.FileName = FileName; return o }
func (o *Options) WithSize(Size string) *Options             { o.Size = Size; return o }
func (o *Options) WithNoAddCaller(NoAddCaller bool) *Options { o.NoAddCaller = NoAddCaller; return o }

func (o *Options) getLevel() zapcore.Level {
	switch o.Level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	case LevelPanic:
		return zapcore.PanicLevel
	}

	return zapcore.DebugLevel
}

func (o *Options) getTimeEncoder() zapcore.TimeEncoder {
	switch strings.ToLower(o.TimeFormat) {
	case "epoch_time_encoder", "epoch_time":
		return zapcore.EpochTimeEncoder
	case "epoch_millis_time_encoder", "epoch_millis_time":
		return zapcore.EpochMillisTimeEncoder
	case "epoch_nanos_time_encoder", "epoch_nanos_time":
		return zapcore.EpochNanosTimeEncoder
	case "iso8601", "":
		return zapcore.ISO8601TimeEncoder
	case "rfc3339":
		return zapcore.RFC3339TimeEncoder
	case "rfc3339nano":
		return zapcore.RFC3339NanoTimeEncoder
	}

	return zapcore.TimeEncoderOfLayout(o.TimeFormat)
}

func (o *Options) getEncoder() zapcore.Encoder {
	enc := zap.NewProductionEncoderConfig()
	enc.EncodeTime = o.getTimeEncoder()
	enc.EncodeLevel = zapcore.CapitalLevelEncoder

	if strings.ToLower(o.Format) == FormatText {
		return zapcore.NewConsoleEncoder(enc)
	}

	return zapcore.NewJSONEncoder(enc)
}

func (o *Options) getWriteSyncer() zapcore.WriteSyncer {
	switch strings.ToLower(o.Output) {
	case OutputStderr:
		return zapcore.AddSync(os.Stderr)

	case OutputFile:
		if len(o.FileName) > 0 {
			writeSyncer := &lumberjack.Logger{
				Filename:   o.FileName,
				MaxSize:    FormatSize(o.Size),
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   false,
			}

			return zapcore.AddSync(writeSyncer)
		}

	default:
	}

	return zapcore.AddSync(os.Stdout)
}

func (o *Options) getAddCaller() []zap.Option {
	if !o.NoAddCaller {
		return []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
	}

	return nil
}

// 50m
const defaultSize = 50 * 1024 * 1024

func FormatSize(size string) int {
	sz := int64(1)
	base := 1

	size = strings.TrimSpace(size)
	if len(size) == 0 {
		return defaultSize
	}

	unit := size[len(size)-1]
	if unit >= '0' && unit <= '9' {
		goto end
	}

	size = size[:len(size)-1]
	if len(size) == 0 {
		size = "1"
	}

	switch unit {
	case 'k', 'K':
		base = 1024
	case 'm', 'M':
		base = 1024 * 1024
	case 'g', 'G':
		base = 1024 * 1024 * 1024
	default:
		return defaultSize
	}

end:
	sz, _ = strconv.ParseInt(size, 10, 64)
	if sz == 0 {
		return defaultSize
	}

	return int(sz) * base
}
