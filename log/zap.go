package log

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Instance *zap.Logger

func Init() {
	Instance = NewLog(config.Instance.Log)
}

func NewLog(cfg config.Log) *zap.Logger {
	var logFile string

	if cfg.Dir != "" {
		logDir, err := filepath.Abs(cfg.Dir)
		if err != nil {
			panic("failed to resolve log dir path: " + err.Error())
		}

		if err := os.MkdirAll(logDir, 0o750); err != nil {
			panic("failed to create log directory: " + err.Error())
		}

		logFile = filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	}

	level := func() zapcore.Level {
		switch strings.ToLower(cfg.Level) {
		case "info":
			return zapcore.InfoLevel
		case "warn", "warning":
			return zapcore.WarnLevel
		case "error":
			return zapcore.ErrorLevel
		case "dpanic":
			return zapcore.DPanicLevel
		case "panic":
			return zapcore.PanicLevel
		case "fatal":
			return zapcore.FatalLevel
		default:
			return zapcore.DebugLevel
		}
	}()

	return zap.New(zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		ConsoleSeparator: " ",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(ansi.ColorCode("white+b") + t.Format("2006-01-02 15:04:05") + ansi.Reset)
		},
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			label := l.CapitalString()
			if color, ok := map[zapcore.Level]string{
				zapcore.DebugLevel:  ansi.ColorCode("magenta"),
				zapcore.InfoLevel:   ansi.ColorCode("green"),
				zapcore.WarnLevel:   ansi.ColorCode("yellow+b"),
				zapcore.ErrorLevel:  ansi.ColorCode("red+b"),
				zapcore.DPanicLevel: ansi.ColorCode("cyan+b"),
				zapcore.PanicLevel:  ansi.ColorCode("white+b+h:red"),
				zapcore.FatalLevel:  ansi.ColorCode("white+b:red"),
			}[l]; ok {
				enc.AppendString("[" + color + label + ansi.Reset + "]")
			} else {
				enc.AppendString("[" + ansi.DefaultBG + label + ansi.Reset + "]")
			}
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}), zapcore.AddSync(colorable.NewColorableStdout()), level), zapcore.NewCore(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}), zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}), level)), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
