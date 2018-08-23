package golog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	megaBytes  = 1024 * 1024
	defMaxSize = 100
)

type FileOut struct {
	// Filename is the file to write logs to.  Backup log files will be
	// retained in the same directory.
	Filename string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int

	size int64
	file *os.File
}

// backupName creates a new filename from the given name, inserting a timestamp
// between the filename and the extension, using the local time if requested
// (otherwise UTC).
func backupName(name string) string {
	var t = time.Now()

	dir := filepath.Dir(name)
	ext := filepath.Ext(name)

	filename := filepath.Base(name)
	prefix := filename[:len(filename)-len(ext)]

	timestamp := t.Format("2006_01_02-15_04_05.000")
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

// max returns the maximum size in bytes of log files before rolling.
func (self *FileOut) max() int64 {
	if self.MaxSize == 0 {
		return int64(defMaxSize * megaBytes)
	}
	return int64(self.MaxSize) * int64(megaBytes)
}

// genFilename generates the name of the logfile from the current time.
func (self *FileOut) filename() string {
	if self.Filename != "" {
		return self.Filename
	}

	var name = filepath.Base(os.Args[0]) + "_defname.log"
	return filepath.Join(os.TempDir(), name)
}

// dir returns the directory for the current filename.
func (self *FileOut) dir() string {
	return filepath.Dir(self.filename())
}

// openNew opens a new log file for writing, moving any old log file out of the
// way.  This methods assumes the file has already been closed.
func (self *FileOut) openNew() error {
	err := os.MkdirAll(self.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := self.filename()
	if _, e := os.Stat(name); e == nil {
		// move the existing file
		var newname = backupName(name)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}
	}

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	mode := os.FileMode(0644)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}

	self.file = f
	self.size = 0
	return nil
}

// rotate closes the current file, moves it aside with a timestamp in the name,
// (if it exists), opens a new file with the original filename, and then runs
// post-rotation processing and removal.
func (self *FileOut) rotate() error {
	if err := self.close(); err != nil {
		return err
	}
	if err := self.openNew(); err != nil {
		return err
	}
	return nil
}

// close closes the file if it is open.
func (self *FileOut) close() error {
	if self.file == nil {
		return nil
	}
	err := self.file.Close()
	self.file = nil
	return err
}

// Write to file
func (self *FileOut) Write(level Level, text string) (int, error) {
	var Logs = []string{
		LevelStr(level),
		text,
	}

	var strLogs = strings.Join(Logs, " ")

	var newsize = int64(len(strLogs)) + self.size
	if newsize > self.max() {
		err := self.rotate()
		if err != nil {
			return 0, err
		}
	}

	return io.WriteString(self.file, strLogs)
}

// Close implements io.Closer, and closes the current logfile.
func (self *FileOut) Close() error {
	return self.close()
}

func NewFileOut(logfile string) *FileOut {
	var fileOut = &FileOut{
		Filename: logfile,
		MaxSize:  100 * megaBytes,
	}

	fileOut.openNew()

	return fileOut
}
