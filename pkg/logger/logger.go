/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap logger wrapper that will be used throughout the logger project
type Logger struct {
	*zap.SugaredLogger
}

var (
	// Global logger instance
	log *Logger
)

// Initialize configures the logger at the start of the project
func Initialize(level string, isDevelopment bool) {
	// Set log level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create Zap configuration
	var config zap.Config
	if isDevelopment {
		// colorful, console-oriented logger for the development environment
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// JSON logger configured for production environment
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	config.Level.SetLevel(zapLevel)

	// Build the logger
	zapLogger, err := config.Build()
	if err != nil {
		panic("logger initialization failed: " + err.Error())
	}

	log = &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// NewWithOptions creates a new logger with the given options
func NewWithOptions(opts ...zap.Option) *Logger {
	if log == nil {
		Initialize("info", true)
	}

	zapLogger := log.SugaredLogger.Desugar().WithOptions(opts...)
	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// With creates a child logger with the given fields
func With(args ...interface{}) *Logger {
	if log == nil {
		Initialize("info", true)
	}
	return &Logger{
		SugaredLogger: log.SugaredLogger.With(args...),
	}
}

// New method creates a new logger instance
func Debug(args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Debug(args...)
}

// Debugf method formats and logs debug messages
func Debugf(format string, args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Debugf(format, args...)
}

// New method creates a new logger instance
func Info(args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Info(args...)
}

// Infof method formats and logs info messages
func Infof(format string, args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Infof(format, args...)
}

// Warn method logs messages at the warn level
func Warn(args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Warn(args...)
}

// Warnf method formats and logs warn messages
func Warnf(format string, args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Warnf(format, args...)
}

// Error method logs messages at the error level
func Error(args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Error(args...)
}

// Errorf method formats and logs error messages
func Errorf(format string, args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Errorf(format, args...)
}

// Fatal method logs messages at the fatal level and exits the program
func Fatal(args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Fatal(args...)
}

// Fatalf method formats and logs fatal messages and exits the program
func Fatalf(format string, args ...interface{}) {
	if log == nil {
		Initialize("info", true)
	}
	log.SugaredLogger.Fatalf(format, args...)
}

// Sync flushes the logger buffer
func Sync() error {
	if log == nil {
		return nil
	}
	return log.SugaredLogger.Sync()
}
