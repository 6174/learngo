package log

import (
	"os"
	"path"
	"fmt"
)

var (
	loggers []*Logger
	// GitLogger logger for git
	GitLogger *Logger
)

// NewLogger create a logger
func NewLogger(bufLen int64, mode, config string) {
	logger := newLogger(bufLen)

	isExist := false
	for i, l := range loggers {
		if l.adapter == mode {
			isExist = true
			loggers[i] = logger
		}
	}
	if !isExist {
		loggers = append(loggers, logger)
	}
	if err := logger.SetLogger(mode, config); err != nil {
		Fatal(2, "Fail to set logger (%s): %v", mode, err)
	}
}

// NewGitLogger create a logger for git
// FIXME: use same log level as other loggers.
func NewGitLogger(logPath string) {
	path := path.Dir(logPath)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		Fatal(4, "Fail to create dir %s: %v", path, err)
	}

	GitLogger = newLogger(0)
	GitLogger.SetLogger("file", fmt.Sprintf(`{"level":0,"filename":"%s","rotate":false}`, logPath))
}

// Trace records trace log
func Trace(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Trace(format, v...)
	}
}

// Debug records debug log
func Debug(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Debug(format, v...)
	}
}

// Info records info log
func Info(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}

// Warn records warnning log
func Warn(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Warn(format, v...)
	}
}

// Error records error log
func Error(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Error(skip, format, v...)
	}
}

// Critical records critical log
func Critical(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Critical(skip, format, v...)
	}
}

// Fatal records error log and exit process
func Fatal(skip int, format string, v ...interface{}) {
	Error(skip, format, v...)
	for _, l := range loggers {
		l.Close()
	}
	os.Exit(1)
}

// Close closes all the loggers
func Close() {
	for _, l := range loggers {
		l.Close()
	}
}