package jaglogger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger interface {
	Critical(...any)
	Criticalf(string, ...any)
	Error(...any)
	Errorf(string, ...any)
	Warning(...any)
	Warningf(string, ...any)
	Notice(...any)
	Noticef(string, ...any)
	Info(...any)
	Infof(string, ...any)
	Debug(...any)
	Debugf(string, ...any)
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota + 1
	LogLevelInfo
	LogLevelNotice
	LogLevelWarning
	LogLevelError
	LogLevelCritical
)

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelCritical:
		return "[CRITICAL]"
	case LogLevelError:
		return "[ERROR]"
	case LogLevelWarning:
		return "[WARNING]"
	case LogLevelNotice:
		return "[NOTICE]"
	case LogLevelInfo:
		return "[INFO]"
	case LogLevelDebug:
		return "[DEBUG]"
	default:
		return fmt.Sprintf("invalid log level: %d", ll)
	}
}

type logger struct {
	outputs map[LogLevel]*log.Logger
}

func (l logger) Critical(v ...any) {
	l.log(LogLevelCritical, v...)
}
func (l logger) Criticalf(format string, v ...any) {
	l.logf(LogLevelCritical, format, v...)
}

func (l logger) Error(v ...any) {
	l.log(LogLevelError, v...)
}
func (l logger) Errorf(format string, v ...any) {
	l.logf(LogLevelError, format, v...)
}

func (l logger) Warning(v ...any) {
	l.log(LogLevelWarning, v...)
}
func (l logger) Warningf(format string, v ...any) {
	l.logf(LogLevelWarning, format, v...)
}

func (l logger) Notice(v ...any) {
	l.log(LogLevelNotice, v...)
}
func (l logger) Noticef(format string, v ...any) {
	l.logf(LogLevelNotice, format, v...)
}

func (l logger) Info(v ...any) {
	l.log(LogLevelInfo, v...)
}
func (l logger) Infof(format string, v ...any) {
	l.logf(LogLevelInfo, format, v...)
}

func (l logger) Debug(v ...any) {
	l.log(LogLevelDebug, v...)
}
func (l logger) Debugf(format string, v ...any) {
	l.logf(LogLevelDebug, format, v...)
}

func (l logger) log(level LogLevel, v ...any) {
	if logOutput, ok := l.outputs[level]; ok {
		logOutput.Output(3, fmt.Sprint(v...))

	}
}

func (l logger) logf(level LogLevel, format string, v ...any) {
	if logOutput, ok := l.outputs[level]; ok {
		logOutput.Output(3, fmt.Sprintf(format, v...))
	}
}

func NewLogger(minLevel LogLevel, opts ...Option) Logger {

	//initialize settings with default values
	loggerSettings := settings{
		LogLevelConfigs: map[LogLevel]Config{
			LogLevelCritical: {},
			LogLevelError:    {},
			LogLevelWarning:  {},
			LogLevelNotice:   {},
			LogLevelInfo:     {},
			LogLevelDebug:    {},
		},
		DefaultErrOutputs:    []io.Writer{os.Stderr},
		DefaultNonErrOutputs: []io.Writer{os.Stdout},
		DefaultFlags:         log.Ldate | log.Ltime | log.Llongfile,
	}

	// Apply passed in settings
	for _, opt := range opts {
		opt(&loggerSettings)
	}

	// Go through the log level configs and replace any empty values with default values.
	for logLevel, conf := range loggerSettings.LogLevelConfigs {
		if conf.Flags == 0 {
			conf.Flags = loggerSettings.DefaultFlags
		}
		if conf.Prefix == "" {
			conf.Prefix = logLevel.String()
		}
		if len(conf.Outputs) == 0 && logLevel >= minLevel {
			if logLevel >= LogLevelWarning {
				conf.Outputs = loggerSettings.DefaultErrOutputs
			} else {
				conf.Outputs = loggerSettings.DefaultNonErrOutputs
			}
		}
		loggerSettings.LogLevelConfigs[logLevel] = conf
	}

	return logger{
		outputs: map[LogLevel]*log.Logger{
			LogLevelCritical: log.New(
				io.MultiWriter(loggerSettings.LogLevelConfigs[LogLevelCritical].Outputs...),
				loggerSettings.LogLevelConfigs[LogLevelCritical].Prefix,
				loggerSettings.LogLevelConfigs[LogLevelCritical].Flags,
			),
			LogLevelError: log.New(
				io.MultiWriter(loggerSettings.LogLevelConfigs[LogLevelError].Outputs...),
				loggerSettings.LogLevelConfigs[LogLevelError].Prefix,
				loggerSettings.LogLevelConfigs[LogLevelError].Flags,
			),
			LogLevelWarning: log.New(
				io.MultiWriter(loggerSettings.LogLevelConfigs[LogLevelWarning].Outputs...),
				loggerSettings.LogLevelConfigs[LogLevelWarning].Prefix,
				loggerSettings.LogLevelConfigs[LogLevelWarning].Flags,
			),
			LogLevelNotice: log.New(
				io.MultiWriter(loggerSettings.LogLevelConfigs[LogLevelNotice].Outputs...),
				loggerSettings.LogLevelConfigs[LogLevelNotice].Prefix,
				loggerSettings.LogLevelConfigs[LogLevelNotice].Flags,
			),
			LogLevelInfo: log.New(
				io.MultiWriter(loggerSettings.LogLevelConfigs[LogLevelInfo].Outputs...),
				loggerSettings.LogLevelConfigs[LogLevelInfo].Prefix,
				loggerSettings.LogLevelConfigs[LogLevelInfo].Flags,
			),
			LogLevelDebug: log.New(
				io.MultiWriter(loggerSettings.LogLevelConfigs[LogLevelDebug].Outputs...),
				loggerSettings.LogLevelConfigs[LogLevelDebug].Prefix,
				loggerSettings.LogLevelConfigs[LogLevelDebug].Flags,
			),
		},
	}
}
