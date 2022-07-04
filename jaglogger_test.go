package jaglogger

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		minLevel LogLevel
		opts     []Option
	}

	testLogFile := &os.File{}

	defaultFlag := log.Ldate | log.Ltime | log.Llongfile
	errOutputs := io.MultiWriter(os.Stderr)
	nonErrOutputs := io.MultiWriter(os.Stdout)
	noOutputs := io.MultiWriter()

	tests := []struct {
		name string
		args args
		want Logger
	}{
		{
			name: "Min Level Debug",
			args: args{
				minLevel: LogLevelDebug,
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(errOutputs, "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(errOutputs, "[ERROR]", defaultFlag),
					LogLevelWarning:  log.New(errOutputs, "[WARNING]", defaultFlag),
					LogLevelNotice:   log.New(nonErrOutputs, "[NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(nonErrOutputs, "[INFO]", defaultFlag),
					LogLevelDebug:    log.New(nonErrOutputs, "[DEBUG]", defaultFlag),
				},
			},
		},
		{
			name: "Min Level Info",
			args: args{
				minLevel: LogLevelInfo,
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(errOutputs, "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(errOutputs, "[ERROR]", defaultFlag),
					LogLevelWarning:  log.New(errOutputs, "[WARNING]", defaultFlag),
					LogLevelNotice:   log.New(nonErrOutputs, "[NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(nonErrOutputs, "[INFO]", defaultFlag),
					LogLevelDebug:    log.New(noOutputs, "[DEBUG]", defaultFlag),
				},
			},
		},
		{
			name: "Min Level Notice",
			args: args{
				minLevel: LogLevelNotice,
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(errOutputs, "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(errOutputs, "[ERROR]", defaultFlag),
					LogLevelWarning:  log.New(errOutputs, "[WARNING]", defaultFlag),
					LogLevelNotice:   log.New(nonErrOutputs, "[NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(noOutputs, "[INFO]", defaultFlag),
					LogLevelDebug:    log.New(noOutputs, "[DEBUG]", defaultFlag),
				},
			},
		},
		{
			name: "Min Level Warning",
			args: args{
				minLevel: LogLevelWarning,
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(errOutputs, "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(errOutputs, "[ERROR]", defaultFlag),
					LogLevelWarning:  log.New(errOutputs, "[WARNING]", defaultFlag),
					LogLevelNotice:   log.New(noOutputs, "[NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(noOutputs, "[INFO]", defaultFlag),
					LogLevelDebug:    log.New(noOutputs, "[DEBUG]", defaultFlag),
				},
			},
		},
		{
			name: "Min Level Error",
			args: args{
				minLevel: LogLevelError,
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(errOutputs, "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(errOutputs, "[ERROR]", defaultFlag),
					LogLevelWarning:  log.New(noOutputs, "[WARNING]", defaultFlag),
					LogLevelNotice:   log.New(noOutputs, "[NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(noOutputs, "[INFO]", defaultFlag),
					LogLevelDebug:    log.New(noOutputs, "[DEBUG]", defaultFlag),
				},
			},
		},
		{
			name: "Min Level Critical",
			args: args{
				minLevel: LogLevelCritical,
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(errOutputs, "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(noOutputs, "[ERROR]", defaultFlag),
					LogLevelWarning:  log.New(noOutputs, "[WARNING]", defaultFlag),
					LogLevelNotice:   log.New(noOutputs, "[NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(noOutputs, "[INFO]", defaultFlag),
					LogLevelDebug:    log.New(noOutputs, "[DEBUG]", defaultFlag),
				},
			},
		},
		{
			name: "Log Level Config Options",
			args: args{
				minLevel: LogLevelDebug,
				opts: []Option{
					SetCriticalLoggerOpt(Config{Outputs: []io.Writer{testLogFile}}),
					SetErrorLoggerOpt(Config{Prefix: "[TEST_ERROR]"}),
					SetWarningLoggerOpt(Config{Flags: log.LstdFlags}),
					SetNoticeLoggerOpt(Config{Outputs: []io.Writer{ioutil.Discard}, Prefix: "[TEST_NOTICE]"}),
					SetInfoLoggerOpt(Config{Prefix: "[TEST_INFO]", Flags: log.LstdFlags}),
					SetDebugLoggerOpt(Config{Outputs: []io.Writer{ioutil.Discard}, Prefix: "[TEST_DEBUG]", Flags: log.LstdFlags}),
				},
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(io.MultiWriter(testLogFile), "[CRITICAL]", defaultFlag),
					LogLevelError:    log.New(errOutputs, "[TEST_ERROR]", defaultFlag),
					LogLevelWarning:  log.New(errOutputs, "[WARNING]", log.LstdFlags),
					LogLevelNotice:   log.New(io.MultiWriter(ioutil.Discard), "[TEST_NOTICE]", defaultFlag),
					LogLevelInfo:     log.New(nonErrOutputs, "[TEST_INFO]", log.LstdFlags),
					LogLevelDebug:    log.New(io.MultiWriter(ioutil.Discard), "[TEST_DEBUG]", log.LstdFlags),
				},
			},
		},
		{
			name: "Set Defaults Options",
			args: args{
				minLevel: LogLevelDebug,
				opts: []Option{
					SetDefaultErrorOutputsOpt([]io.Writer{testLogFile}),
					SetDefaultNonErrorOutputOpt([]io.Writer{ioutil.Discard}),
					SetDefaultFlagsOpt(log.LstdFlags),
				},
			},
			want: logger{
				outputs: map[LogLevel]*log.Logger{
					LogLevelCritical: log.New(io.MultiWriter(testLogFile), "[CRITICAL]", log.LstdFlags),
					LogLevelError:    log.New(io.MultiWriter(testLogFile), "[ERROR]", log.LstdFlags),
					LogLevelWarning:  log.New(io.MultiWriter(testLogFile), "[WARNING]", log.LstdFlags),
					LogLevelNotice:   log.New(io.MultiWriter(ioutil.Discard), "[NOTICE]", log.LstdFlags),
					LogLevelInfo:     log.New(io.MultiWriter(ioutil.Discard), "[INFO]", log.LstdFlags),
					LogLevelDebug:    log.New(io.MultiWriter(ioutil.Discard), "[DEBUG]", log.LstdFlags),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogger(tt.args.minLevel, tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_logger_Critical(t *testing.T) {
	type args struct {
		v []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetCriticalLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[CRITICAL\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Critical(tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Criticalf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetCriticalLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{format: "test %s", v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[CRITICAL\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Criticalf(tt.args.format, tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Error(t *testing.T) {
	type args struct {
		v []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetErrorLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[ERROR\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Error(tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Errorf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetErrorLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{format: "test %s", v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[ERROR\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Errorf(tt.args.format, tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Warning(t *testing.T) {
	type args struct {
		v []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetWarningLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[WARNING\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Warning(tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Warningf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetWarningLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{format: "test %s", v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[WARNING\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Warningf(tt.args.format, tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Notice(t *testing.T) {
	type args struct {
		v []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetNoticeLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[NOTICE\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Notice(tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Noticef(t *testing.T) {
	type args struct {
		format string
		v      []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetNoticeLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{format: "test %s", v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[NOTICE\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Noticef(tt.args.format, tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Info(t *testing.T) {
	type args struct {
		v []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetInfoLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[INFO\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Info(tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Infof(t *testing.T) {
	type args struct {
		format string
		v      []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetInfoLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{format: "test %s", v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[INFO\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Infof(tt.args.format, tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Debug(t *testing.T) {
	type args struct {
		v []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetDebugLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[DEBUG\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debug(tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func Test_logger_Debugf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}

	loggerOutput := new(bytes.Buffer)
	tests := []struct {
		name      string
		l         Logger
		args      args
		wantMatch *regexp.Regexp
	}{
		{
			name:      "Test Output",
			l:         NewLogger(LogLevelCritical, SetDebugLoggerOpt(Config{Outputs: []io.Writer{loggerOutput}})),
			args:      args{format: "test %s", v: []any{"test"}},
			wantMatch: regexp.MustCompile(`^\[DEBUG\]\d{4}\/\d{2}\/\d{2}\s\d{2}\:\d{2}\:\d{2}\s.*\/jaglogger_test\.go\:\d*\:\stest\stest\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debugf(tt.args.format, tt.args.v...)
			got := loggerOutput.String()
			assert.Regexp(t, tt.wantMatch, got)
		})
	}
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name string
		ll   LogLevel
		want string
	}{
		{
			name: "Critical Level",
			ll:   LogLevelCritical,
			want: "[CRITICAL]",
		},
		{
			name: "Error Level",
			ll:   LogLevelError,
			want: "[ERROR]",
		},
		{
			name: "Warning Level",
			ll:   LogLevelWarning,
			want: "[WARNING]",
		},
		{
			name: "Notice Level",
			ll:   LogLevelNotice,
			want: "[NOTICE]",
		},
		{
			name: "Info Level",
			ll:   LogLevelInfo,
			want: "[INFO]",
		},
		{
			name: "Debug Level",
			ll:   LogLevelDebug,
			want: "[DEBUG]",
		},
		{
			name: "Bad Level Value",
			ll:   LogLevel(99),
			want: "invalid log level: 99",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ll.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
