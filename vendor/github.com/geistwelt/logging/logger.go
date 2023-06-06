package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-stack/stack"
)

type Logger interface {
	Debug(msg string, kvs ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg string, kvs ...interface{})
	Infof(format string, args ...interface{})
	Warn(msg string, kvs ...interface{})
	Warnf(format string, args ...interface{})
	Error(msg string, kvs ...interface{})
	Errorf(format string, args ...interface{})
	Panic(msg string, kvs ...interface{})
	Panicf(format string, args ...interface{})

	SetModule(module string, level LogLevel)
	DeriveChildLogger(module string) Logger
	Update(opt Option) error
}

type Option struct {
	Module         string
	FilterLevel    LogLevel
	Spec           string
	FormatSelector string
	Writer         io.Writer
}

type logger struct {
	mutex          sync.RWMutex
	module         string
	filterLevel    LogLevel
	modules        map[string]LogLevel
	formatter      Formatter
	formatSelector string // 选择以什么格式输出日志，json 或者 terminal。
	writer         io.Writer
}

func MustNewLogger(opts ...Option) Logger {
	l, err := NewLogger(opts...)
	if err != nil {
		panic(err)
	}
	return l
}

func NewLogger(opts ...Option) (Logger, error) {
	l := &logger{
		filterLevel:    InfoLevel,
		modules:        make(map[string]LogLevel),
		formatSelector: "terminal",
		writer:         os.Stdout,
	}
	formatters, err := ParseFormat(defaultFormat)
		if err != nil {
			return nil, err
		}
		l.formatter = NewMultiFormatter(formatters...)
	if len(opts) > 0 {
		l.Update(opts[0])	
	}
	return l, nil
}

func (l *logger) SetModule(module string, level LogLevel) {
	l.mutex.Lock()
	if l.module == module {
		l.filterLevel = level
	}
	l.modules[module] = level
	l.mutex.Unlock()
}

func (l *logger) DeriveChildLogger(module string) Logger {
	child := l.clone()
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if module == l.module {
		return child
	}
	child.module = module
	_, ok := l.modules[module]
	if ok {
		child.filterLevel = l.modules[module]
	} else {
		child.filterLevel = l.filterLevel
	}
	return child
}

func (l *logger) Update(opt Option) error {
	if opt.Module != "" {
		l.module = opt.Module
	}
	if opt.FilterLevel != 0 {
		l.filterLevel = opt.FilterLevel
		if l.module != "" {
			l.modules[l.module] = l.filterLevel
		}
	}
	if opt.Spec != "" {
		formatters, err := ParseFormat(opt.Spec)
		if err != nil {
			return err
		}
		l.formatter = NewMultiFormatter(formatters...)
	}
	if opt.FormatSelector != "" {
		switch opt.FormatSelector {
		case "terminal", "":
			l.formatSelector = "terminal"
		case "json":
			l.formatSelector = "json"
		default:
			return fmt.Errorf("unknown format selector: (%s)", opt.FormatSelector)
		}
		l.formatSelector = opt.FormatSelector
	}
	if opt.Writer != nil {
		l.writer = opt.Writer
	}
	return nil
}

func (l *logger) Debug(msg string, kvs ...interface{}) {
	l.mutex.RLock()
	l.log(msg, DebugLevel, kvs)
	l.mutex.RUnlock()
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.mutex.RLock()
	l.log(fmt.Sprintf(format, args...), DebugLevel, nil)
	l.mutex.RUnlock()
}

func (l *logger) Info(msg string, kvs ...interface{}) {
	l.mutex.RLock()
	l.log(msg, InfoLevel, kvs)
	l.mutex.RUnlock()
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.mutex.RLock()
	l.log(fmt.Sprintf(format, args...), InfoLevel, nil)
	l.mutex.RUnlock()
}

func (l *logger) Warn(msg string, kvs ...interface{}) {
	l.mutex.RLock()
	l.log(msg, WarnLevel, kvs)
	l.mutex.RUnlock()
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.mutex.RLock()
	l.log(fmt.Sprintf(format, args...), WarnLevel, nil)
	l.mutex.RUnlock()
}

func (l *logger) Error(msg string, kvs ...interface{}) {
	l.mutex.RLock()
	l.log(msg, ErrorLevel, kvs)
	l.mutex.RUnlock()
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.mutex.RLock()
	l.log(fmt.Sprintf(format, args...), ErrorLevel, nil)
	l.mutex.RUnlock()
}

func (l *logger) Panic(msg string, kvs ...interface{}) {
	l.mutex.RLock()
	l.log(msg, PanicLevel, kvs)
	l.mutex.RUnlock()
	panic(msg)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.mutex.RLock()
	l.log(fmt.Sprintf(format, args...), PanicLevel, nil)
	l.mutex.RUnlock()
	panic(fmt.Sprintf(format, args...))
}

func (l *logger) log(msg string, level LogLevel, keyValues []interface{}) {
	if level > l.filterLevel {
		return
	}
	entry := Entry{
		Level:          level,
		Time:           time.Now(),
		Module:         l.module,
		Message:        msg,
		Call:           stack.Caller(2),
		FormatSelector: l.formatSelector,
		KeyValues:      keyValues,
	}
	l.formatter.Format(l.writer, entry)
}

func (l *logger) clone() *logger {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	cpy := &logger{
		module:         l.module,
		filterLevel:    l.filterLevel,
		modules:        make(map[string]LogLevel),
		formatter:      l.formatter.(*MultiFormatter).clone(),
		formatSelector: l.formatSelector,
		writer:         l.writer,
	}
	for module, level := range l.modules {
		cpy.modules[module] = level
	}
	return cpy
}
