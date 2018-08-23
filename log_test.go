package golog

import (
	"log"
	"testing"
)

func TestLevel(t *testing.T) {
	logex := New("test")
	SetFileLogByString("test", "./")

	SetStackByString("test", true)
	logex.Debugf("%d %s %v\n", 1, "hello", map[int]int{1: 34, 2: 48})
	logex.Errorln("hello1")
	logex.Infoln("no")

	SetStackByString("test", false)
	logex.Errorln("2")
	logex.Warnln("warn", 1, 4.5)
	logex.Fatalf("%s %s\n", "sss", "kkk")
}

func TestMyLog(t *testing.T) {
	logex := New("testlog")
	logex.Debugln("hello1")
	logex.Debugln("hello2")
}

func TestSystemLog(t *testing.T) {
	log.Println("hello1")
	log.Println("hello2")
	log.Println("hello3")
}
