package jaglogger

import "io"

// Config holds the data that will be used to build the logger of a specfic log level.
type Config struct {
	Outputs []io.Writer
	Prefix  string
	Flags   int
}

// Option is a function type that allows modifications of settings for the logger
type Option func(*settings)

// settings holds the properties that can be modified by the Option type
type settings struct {
	LogLevelConfigs      map[LogLevel]Config
	DefaultErrOutputs    []io.Writer
	DefaultNonErrOutputs []io.Writer
	DefaultFlags         int
}

// SetCriticalLoggerOpt sets the logger configuration for the "Critical" log level
func SetCriticalLoggerOpt(config Config) Option {
	return func(s *settings) {
		s.LogLevelConfigs[LogLevelCritical] = config
	}
}

func SetErrorLoggerOpt(config Config) Option {
	return func(s *settings) {
		s.LogLevelConfigs[LogLevelError] = config
	}
}

func SetWarningLoggerOpt(config Config) Option {
	return func(s *settings) {
		s.LogLevelConfigs[LogLevelWarning] = config
	}
}

func SetNoticeLoggerOpt(config Config) Option {
	return func(s *settings) {
		s.LogLevelConfigs[LogLevelNotice] = config
	}
}

func SetInfoLoggerOpt(config Config) Option {
	return func(s *settings) {
		s.LogLevelConfigs[LogLevelInfo] = config
	}
}

func SetDebugLoggerOpt(config Config) Option {
	return func(s *settings) {
		s.LogLevelConfigs[LogLevelDebug] = config
	}
}

func SetDefaultErrorOutputsOpt(outputs []io.Writer) Option {
	return func(s *settings) {
		s.DefaultErrOutputs = outputs
	}
}

func SetDefaultNonErrorOutputOpt(outputs []io.Writer) Option {
	return func(s *settings) {
		s.DefaultNonErrOutputs = outputs
	}
}

func SetDefaultFlagsOpt(flag int) Option {
	return func(s *settings) {
		s.DefaultFlags = flag
	}
}
