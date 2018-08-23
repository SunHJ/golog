package golog

import (
	"errors"
	"path/filepath"
	"sync"
)

const (
	logExt = ".log"
)

var (
	logMap      = map[string]*Logger{}
	logMapGuard sync.RWMutex

	fileSet      = map[string]*FileOut{}
	fileSetGuard sync.RWMutex
)

func add(l *Logger) {
	logMapGuard.Lock()
	defer logMapGuard.Unlock()

	if _, ok := logMap[l.name]; ok {
		panic("duplicate logger name:" + l.name)
	}

	logMap[l.name] = l
}

func regFileLog(logdir string, name string) *FileOut {
	fileSetGuard.Lock()
	defer fileSetGuard.Unlock()

	var logFile *FileOut

	var logname = filepath.Join(logdir, name+logExt)
	if l, ok := fileSet[name]; ok {
		l.Filename = logname
		logFile = l
	} else {
		logFile = NewFileOut(logname)
		fileSet[name] = logFile
	}

	return logFile
}

func str2loglevel(level string) Level {
	switch level {
	case "debug":
		return Level_Debug
	case "info":
		return Level_Info
	case "warn":
		return Level_Warn
	case "error":
		return Level_Error
	case "fatal":
		return Level_Fatal
	}

	return Level_Debug
}

func VisitLogger(name string, callback func(*Logger) bool) error {
	logMapGuard.RLock()
	defer logMapGuard.RUnlock()

	if name == "*" {
		for _, l := range logMap {
			callback(l)
		}
	} else {
		if l, ok := logMap[name]; ok {
			callback(l)
		} else {
			return errors.New("logger not found")
		}
	}

	return nil
}

// 通过字符串设置某一类日志的堆栈depth
func SetCallDepthByString(loggerName string, depth int) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetCallDepth(depth)
		return true
	})
}

// 通过字符串设置某一类日志的级别
func SetLevelByString(loggerName string, level string) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetLevelByString(level)
		return true
	})
}

// 通过字符串设置stackinfo
func SetStackByString(loggerName string, showStack bool) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetStackFlag(showStack)
		return true
	})
}

// 通过字符串设置filelog
func SetFileLogByString(loggerName string, logPath string) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetFileFlag(logPath)
		return true
	})
}

func SetConsoleLogByString(loggerName string, enable bool) error {

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetConsoleFlag(enable)
		return true
	})
}

func EnableConsole(enable bool) {
	gConsole.Disable = !enable
}

func ClearAll() {
	logMapGuard.Lock()
	logMap = map[string]*Logger{}
	logMapGuard.Unlock()
}
