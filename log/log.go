package log

import (
	"log"
	"sync"
)

type Logger interface {
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	SetLevel(level Level)
}

type Level int

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var gmu sync.Mutex
var gl Logger = &DefaultLogger{level: InfoLevel}

func SetLogger(logger Logger) {
	gmu.Lock()
	gl = logger
	gmu.Unlock()
}

func Panicf(format string, v ...interface{}) {
	gl.Panicf(format, v...)
}
func Fatalf(format string, v ...interface{}) {
	gl.Fatalf(format, v...)

}
func Errorf(format string, v ...interface{}) {
	gl.Errorf(format, v...)
}
func Infof(format string, v ...interface{}) {
	gl.Infof(format, v...)
}
func Warnf(format string, v ...interface{}) {
	gl.Warnf(format, v...)
}
func Debugf(format string, v ...interface{}) {
	gl.Debugf(format, v...)
}

type DefaultLogger struct {
	level Level
	mu    sync.Mutex
}

func (l *DefaultLogger) printf(format string, v []interface{}) {
	log.Printf(format, v...)
}

func (l *DefaultLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *DefaultLogger) Panicf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= PanicLevel {
		l.printf(format, v)
	}
}
func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= FatalLevel {
		l.printf(format, v)
	}

}
func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= ErrorLevel {
		l.printf(format, v)
	}
}
func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= InfoLevel {
		l.printf(format, v)
	}
}
func (l *DefaultLogger) Warnf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= InfoLevel {
		l.printf(format, v)
	}
}
func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= DebugLevel {
		l.printf(format, v)
	}
}
