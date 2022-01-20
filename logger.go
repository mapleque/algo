package algo

import (
	"fmt"
	"time"
)

const (
	green  = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	blue   = "\033[97;44m"
	reset  = "\033[0m"
)

type LogLevel int

var LogPrefix = map[LogLevel]string{
	Debug: blue + "Debug" + reset,
	Trace: white + "Trace" + reset,
	Info:  green + "Info " + reset,
	Warn:  yellow + "Warn " + reset,
	Error: red + "Error" + reset,
}

const (
	Debug LogLevel = iota
	Trace
	Info
	Warn
	Error
)

type Log struct {
	level LogLevel
}

var log *Log = &Log{level: Trace}

func SetLogLevel(level LogLevel) {
	log.SetLogLevel(level)
}

func (log *Log) SetLogLevel(level LogLevel) {
	log.level = level
}

func (log *Log) Debug(msg ...interface{}) {
	log.Log(Debug, msg...)
}

func (log *Log) Trace(msg ...interface{}) {
	log.Log(Trace, msg...)
}

func (log *Log) Info(msg ...interface{}) {
	log.Log(Info, msg...)
}

func (log *Log) Warn(msg ...interface{}) {
	log.Log(Warn, msg...)
}

func (log *Log) Error(msg ...interface{}) {
	log.Log(Error, msg...)
}

func (log *Log) Log(level LogLevel, msg ...interface{}) {
	if log.level > level {
		return
	}
	fmt.Printf("%s[%s] %s\n", now(), LogPrefix[level], fmt.Sprint(msg...))
}

func (log *Log) Debugf(fmt string, msg ...interface{}) {
	log.Logf(Debug, fmt, msg...)
}

func (log *Log) Tracef(fmt string, msg ...interface{}) {
	log.Logf(Trace, fmt, msg...)
}

func (log *Log) Infof(fmt string, msg ...interface{}) {
	log.Logf(Info, fmt, msg...)
}

func (log *Log) Warnf(fmt string, msg ...interface{}) {
	log.Logf(Warn, fmt, msg...)
}

func (log *Log) Errorf(fmt string, msg ...interface{}) {
	log.Logf(Error, fmt, msg...)
}
func (log *Log) Logf(level LogLevel, format string, msg ...interface{}) {
	if log.level > level {
		return
	}
	fmt.Printf("%s[%s] %s\n", now(), LogPrefix[level], fmt.Sprintf(format, msg...))
}

func now() string {
	return time.Now().Format("150405.000")
}
