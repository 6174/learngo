package log

import (
	"sync"
	"fmt"
	"runtime"
	"strings"
	"path/filepath"
	"os"
)

// .___        __                 _____
// |   | _____/  |_  ____________/ ____\____    ____  ____
// |   |/    \   __\/ __ \_  __ \   __\\__  \ _/ ___\/ __ \
// |   |   |  \  | \  ___/|  | \/|  |   / __ \\  \__\  ___/
// |___|___|  /__|  \___  >__|   |__|  (____  /\___  >___  >
//          \/          \/                  \/     \/    \/
var adapters = make(map[string]loggerType)

// Register registers given logger provider to adapters.
func Register(name string, log loggerType) {
	if log == nil {
		panic("log: register provider is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("log: register called twice for provider \"" + name + "\"")
	}
	adapters[name] = log
}

// log levels
const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
	FATAL
)

type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, skip, level int) error
	Destroy()
	Flush()
}

type loggerType func() LoggerInterface

type logMsg struct {
	skip  int
	level int
	msg   string
}

type Logger struct {
	adapter  string
	lock     sync.Mutex
	level    int
	msg      chan *logMsg
	outputs  map[string]LoggerInterface
	quit     chan bool
}

func newLogger(buffer int64) *Logger {
	l := &Logger {
		msg:      make(chan *logMsg, buffer),
		outputs:  make(map[string]LoggerInterface),
		quit:     make(chan bool),
	}
	go l.StartLogger()
	return l
}

func (l *Logger) SetLogger(adapter string, config string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if log, ok := adapters[adapter]; ok {
		lg := log()
		if err := lg.Init(config); err != nil {
			return err
		}
		l.outputs[adapter] = lg
		l.adapter = adapter
	} else {
		panic("log: unknow adapter " + adapter)
	}
	return nil
}

func (l *Logger) writerMsg(skip, level int, msg string) error {
	if l.level > level {
		return nil
	}
	lm := &logMsg {
		skip: skip,
		level: level,
	}

	if lm.level >= ERROR {
		pc, file, line, ok := runtime.Caller(skip)
		if ok {
			fn := runtime.FuncForPC(pc)
			var fnName string
			if fn == nil {
				fnName = "?()"
			} else {
				fnName = strings.TrimLeft(filepath.Ext(fn.Name()), ".") + "()"
			}

			fileName := file

			lm.msg = fmt.Sprintf("[%s:%d %s] %s", fileName, line, fnName, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	l.msg <- lm
	return nil
}

// DelLogger removes a logger adapter instance.
func (l *Logger) DelLogger(adapter string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if lg, ok := l.outputs[adapter]; ok {
		lg.Destroy()
		delete(l.outputs, adapter)
	} else {
		panic("log: unknown adapter \"" + adapter + "\" (forgotten register?)")
	}
	return nil
}

// After Logger started , logger will watch msg channel
//  when logMsg come in
//    outputs all logMsg
func (l *Logger) StartLogger() {
	for {
		select {
		case bm := <- l.msg:
			for _, out := range l.outputs {
				if err := out.WriteMsg(bm.msg, bm.skip, bm.level); err != nil {
					fmt.Println("Error, unable to WriteMsg", err)
				}
			}
		case <- l.quit:
			return
		}
	}
}

func (l *Logger) Flush() {
	for _, out := range l.outputs {
		out.Flush()
	}
}

func (l *Logger) Close() {
	l.quit <- true
	// log rest msg
	for {
		if len(l.msg) > 0 {
			bm := <-l.msg
			for _, out := range l.outputs {
				if err := out.WriteMsg(bm.msg, bm.skip, bm.level); err != nil {
					fmt.Println("Error, unable to WriteMsg:", err)
				}
			}
		} else {
			break
		}
	}
	for _, out := range l.outputs {
		out.Flush()
		out.Destroy()
	}
}

// Trace records trace log
func (l *Logger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	l.writerMsg(0, TRACE, msg)
}

// Debug records debug log
func (l *Logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	l.writerMsg(0, DEBUG, msg)
}

// Info records information log
func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	l.writerMsg(0, INFO, msg)
}

// Warn records warnning log
func (l *Logger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	l.writerMsg(0, WARN, msg)
}

// Error records error log
func (l *Logger) Error(skip int, format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	l.writerMsg(skip, ERROR, msg)
}

// Critical records critical log
func (l *Logger) Critical(skip int, format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	l.writerMsg(skip, CRITICAL, msg)
}

// Fatal records error log and exit the process
func (l *Logger) Fatal(skip int, format string, v ...interface{}) {
	msg := fmt.Sprintf("[F] "+format, v...)
	l.writerMsg(skip, FATAL, msg)
	l.Close()
	os.Exit(1)
}