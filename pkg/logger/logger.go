package logger

import (
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level int

func (l Level) String() string {
	switch l {
	case Info:
		return "INFO"
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Fatal:
		return "FATAL"
	default:
		return "UNDEFINED"
	}
}

const (
	Info Level = iota
	Warn
	Error
	Fatal
)

type Logger struct {
	out   io.Writer
	level Level
	mu    sync.Mutex
}

func New(out io.Writer, level Level) *Logger {
	return &Logger{
		out:   out,
		level: level,
	}
}

func (l *Logger) Info(message string, meta ...any) (int, error) {
	return l.print(Info, message, meta...)
}

func (l *Logger) Warn(message string, meta ...any) (int, error) {
	return l.print(Warn, message, meta...)
}

func (l *Logger) Error(err error, meta ...any) (int, error) {
	return l.print(Error, err.Error(), meta...)
}

func (l *Logger) Fatal(err error, meta ...any) (int, error) {
	return l.print(Fatal, err.Error(), meta...)
}

func (l *Logger) print(level Level, message string, meta ...any) (int, error) {
	if level < l.level {
		return 0, nil
	}

	t := time.Now().UTC().Format(time.RFC3339)

	var line string

	if level == Error || level == Fatal {
		trace := string(debug.Stack())
		line = fmt.Sprintf("%s %s: %s %v \n\n %s\n", t, level.String(), message, meta, trace)
	} else {
		line = fmt.Sprintf("%s %s: %s %v\n", t, level.String(), message, meta)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write([]byte(line))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(Error, string(message))
}
