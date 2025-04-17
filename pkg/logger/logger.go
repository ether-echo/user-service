package logger

import (
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
)

var (
	logger      *zap.Logger
	once        sync.Once
	atomicLevel zap.AtomicLevel
)

func BuildLogger(logLevel string) {
	once.Do(func() {
		atomicLevel = zap.NewAtomicLevel()
		SetLevel(logLevel)
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time" // Ключ для времени в JSON
		encoderCfg.EncodeTime = CustomTimeEncoder
		logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, atomicLevel), zap.AddCaller())
	})
}

func SetLevel(logLevel string) {
	switch strings.ToUpper(logLevel) {
	case LevelDebug:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case LevelInfo:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	default:
		panic("invalid log level specified for logger")
	}
}

func CurrentLevel() string {
	return atomicLevel.String()
}

// Logger returns a global logger defined in this package.
// If logger is nil function returns a logger with DEBUG level.
func Logger() *zap.Logger {
	if logger == nil {
		BuildLogger(LevelDebug)
	}

	return logger
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
