package golog

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogFlag uint8

const (
	Ldate     LogFlag         = 1 << iota // the date: 2006/01/02
	Ltime                                 // the time: 15:04:05.000
	Lstack                                // file name and line number: /a/b/c/d.go:23
	Lconsole                              // write to console
	Lfile                                 // write to file
	LstdFlags = Ldate | Ltime             // initial values for the standard logger
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	name   string
	flag   LogFlag // properties
	level  Level
	depth  int
	preidx int

	mu sync.Mutex // ensures atomic writes; protects the following fields
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(name string) *Logger {
	l := &Logger{
		name:  name,
		flag:  LstdFlags,
		level: Level_Debug,
		depth: 3,
	}

	add(l)

	return l
}

func (self *Logger) Output(level Level, text string) {
	now := time.Now() // get this early

	self.mu.Lock()
	defer self.mu.Unlock()

	var strLog string

	if self.flag&Ldate != 0 {
		strLog += now.Format("2006-01-02")
	}

	if self.flag&Ltime != 0 {
		strLog += now.Format(" 15:04:05.000")
	}

	if self.flag&Lstack != 0 {
		// release lock while getting caller info - it'text expensive.
		self.mu.Unlock()

		var file string
		var line int
		var ok bool
		_, file, line, ok = runtime.Caller(self.depth)
		if !ok {
			file = "???"
			line = 0
		}
		strLog += fmt.Sprintf(" %s:%d", file[self.preidx:], line)

		self.mu.Lock()
	}

	strLog += " " + text

	// console log
	if self.flag&Lconsole > 0 {
		gConsole.Write(level, self.name+" "+strLog)
	}

	// file log
	if self.flag&Lfile > 0 {
		if logout, ok := fileSet[self.name]; ok {
			logout.Write(level, strLog)
		}
	}
}

func (self *Logger) Log(level Level, format string, v ...interface{}) {
	if level < self.level {
		return
	}

	var text string

	if format == "" {
		text = fmt.Sprintln(v...)
	} else {
		text = fmt.Sprintf(format+"\n", v...)
	}

	self.Output(level, text)

	if level >= Level_Fatal {
		panic(text)
	}
}

func (self *Logger) Debugf(format string, v ...interface{}) {

	self.Log(Level_Debug, format, v...)
}

func (self *Logger) Debugln(v ...interface{}) {

	self.Log(Level_Debug, "", v...)
}

func (self *Logger) Infof(format string, v ...interface{}) {

	self.Log(Level_Info, format, v...)
}

func (self *Logger) Infoln(v ...interface{}) {

	self.Log(Level_Info, "", v...)
}

func (self *Logger) Warnf(format string, v ...interface{}) {

	self.Log(Level_Warn, format, v...)
}

func (self *Logger) Warnln(v ...interface{}) {
	self.Log(Level_Warn, "", v...)
}

func (self *Logger) Errorf(format string, v ...interface{}) {

	self.Log(Level_Error, format, v...)
}

func (self *Logger) Errorln(v ...interface{}) {
	self.Log(Level_Error, "", v...)
}

func (self *Logger) Fatalf(format string, v ...interface{}) {

	self.Log(Level_Fatal, format, v...)
}

func (self *Logger) Fatalln(v ...interface{}) {

	self.Log(Level_Fatal, "", v...)
}

func (self *Logger) SetCallDepth(depth int) {
	self.depth = depth
}

func (self *Logger) SetStackFlag(stack bool) {
	if stack {
		self.flag |= Lstack
		if _, file, _, ok := runtime.Caller(2); ok {
			self.preidx = strings.Index(file, "src")
			if self.preidx < 0 {
				self.preidx = 0
			}
		}
	} else {
		self.flag &= (^Lstack)
	}
}

func (self *Logger) SetFileFlag(logpath string) {
	if len(logpath) > 0 {
		regFileLog(logpath, self.name)
		self.flag |= Lfile
	} else {
		self.flag &= (^Lfile)
		if l, ok := fileSet[self.name]; ok {
			l.Close()
			delete(fileSet, self.name)
		}
	}
}

func (self *Logger) SetConsoleFlag(console bool) {
	if console {
		self.flag |= Lconsole
	} else {
		self.flag &= (^Lconsole)
	}
}

func (self *Logger) SetLevelByString(level string) {
	level = strings.ToLower(level)
	self.SetLevel(str2loglevel(level))
}

func (self *Logger) SetLevel(lv Level) {
	self.level = lv
}

func (self *Logger) Level() Level {
	return self.level
}
