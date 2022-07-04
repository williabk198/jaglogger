[![go.mod](https://img.shields.io/github/go-mod/go-version/williabk198/jaglogger)](go.mod)
[![go report](https://goreportcard.com/badge/github.com/williabk198/jaglogger)](https://goreportcard.com/report/github.com/williabk198/jaglogger)
[![test status](https://github.com/williabk198/jaglogger/workflows/test/badge.svg)](https://github.com/williabk198/jaglogger/actions)
[![LICENSE](https://img.shields.io/github/license/williabk198/jaglogger)](LICENSE) 

# Just Another Go(JAG) Logger
JAG Logger is built on top of Go's existing `log` package, 
and aims to provide a simple and extensible solution for your logging needs.

## Usage

### Default Logger
If you just want a basic logger to use for your application, you can initialize the JAG Logger like so:
```go
logger := jaglogger.NewLogger(jaglogger.LogLevelInfo)
```
The parameter handed to `NewLogger` is the lowest logger level priority that will be handled. 
So, for the above example, only log priority levels of `Info` and above will be handled. 
The others (just the `Debug` log priority level, in this case) will be discarded.

Using this default impolementation, the log levels of `Critical`, `Error` and `Warning` will print to 
`os.Stderr`.  The `Notice`, `Info` and `Debug` log levels will print to `os.Stdout`. 
Also, using this default implementation, the log entries will be prefixed with `[<level_name>]`
where `<level_name>` is the log level name in all caps (e.g. `Critical` will be `[CRITICAL]`), 
and the log flags being used are `log.Ldate`, `log.Ltime`, `log.Llongfile`.

Meaning the output will look something akin to this:
```
[INFO]2022/07/03 22:05:03 /path/to/workspace/main.go:8: test
```

### Customized Logger
If the default logger that JAG Logger provides isn't quite what yor are looking for,
there are a handful of ways, to customize your experience.

#### Setting Log Level Customizations
To customize what gets printed out and where at each logging level, 
you can utilize the folllowing functions when inializing the logger with `jaglogger.NewLogger`: 
`SetCriticalLoggerOpt`, `SetErrorLoggerOpt`, `SetWarningLoggerOpt`, `SetNoticeLoggerOpt`, `SetInfoLoggerOpt`,
and `SetDebugLoggerOpt`.

**_Example_**:
```go
file, err := os.Create("./error.log")
if err != nil {
  // handle error...
}

logger := jaglogger.NewLogger(
  jaglogger.LogLevelInfo,
  SetErrorLoggerOpt(jaglogger.Config{Outputs: []io.Writer{file}}),
)
```
Along with the output(s) the logger prints to, you can also up date the log flags, and the prefix of the logger using the `Flags` and `Prefix` properties of `jaglogger.Config` struct respectively.

**_Another Example_**:
```go
logger := jaglogger.NewLogger(
  jaglogger.LogLevelInfo,
  SetErrorLoggerOpt(
    jaglogger.Config{Outputs: []io.Writer{file}, Prefix: "SOME_ERROR_PREFIX", Flags: log.LstdFlags},
  ),
)
```

**_NOTE_**: Any values left blank in the `jaglogger.Config` struct will be filled in with default values during
the execution of `jaglogger.NewLogger`

#### Writing logs to multiple locations
As yoy may have noticed with the previous examples, you can have logs write to more than one place if need be. If the need arises in which you do need to write a log to more than one place, this is how you can do so:
```go
file, err := os.Create("./error.log")
if err != nil {
  // handle error...
}

jaglogger.NewLogger(
  jaglogger.LogLevelInfo,
  SetErrorLoggerOpt(
    jaglogger.Config{Outputs: []io.Writer{file, os.Stderr}},
  ),
)
```
Yep, it's really that simple. As long as where you need to write to implements the `io.Writer` interface,
you can just simply plug it into the `Outputs` property of the `jaglogger.Config`.


#### Overriding Defaults
If you wish to update the default values that are used when a `jaglogger.Config` field is left balnk, 
then JAG Logger has you covered there as well. These are the following functions you can pass to the 
`jaglogger.NewLogger` function to modify the default values: `SetDefaultErrorOutputsOpt`,
`SetDefaultNonErrorOutputsOpt`, and `SetDefaultFlagsOpt`

**_Example_**:
```go
logFile, err := os.Create("./some.log")
if err != nil {
  // handle error...
}

errLogFile, err := os.Create("./error.log")
if err != nil {
  // handle error...
}

logger := jaglogger.NewLogger(
  jaglogger.LogLevelInfo,
  jaglogger.SetDefaultNonErrorOutputsOpt([]io.Writer{logFile}),
  jaglogger.SetDefaultErrorOutputsOpt([]io.Witer{errLogFile}),
  jaglogger.SetDefaultFlagsOpt(log.LstdFlags),
)
```

The above will result in the `Critical`, `Error` and `Warning` log levels being printed to the `error.log` file
with the `log.LstdFlags` flags, and `Notice`, `Info` and `Debug` will be writen to the `some.log` file in with
the `log.LstdFlag` as well.
