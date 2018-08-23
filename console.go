package golog

import (
	"io"
	"os"
	"strings"
	"sync"
)

type ConsoleOut struct {
	sync.Mutex

	Disable bool
}

func (c *ConsoleOut) Write(level Level, text string) (int, error) {
	c.Lock()
	defer c.Unlock()

	var Logs = []string{
		LevelColorStr(level),
		text,
	}

	if c.Disable {
		return 0, nil
	}

	return io.WriteString(os.Stdout, strings.Join(Logs, " "))
}

var gConsole = &ConsoleOut{}
