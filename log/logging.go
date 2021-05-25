package log

import (
	l "log"
	"os"
	"path/filepath"

	"github.com/rchargel/hdfs-explorer/base"
)

var (
	Info  *l.Logger
	Warn  *l.Logger
	Error *l.Logger
)

func init() {
	logpath := filepath.Join(base.HomeDir, "log.txt")
	file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		l.Fatal(err)
	}

	Info = l.New(file, "[ INFO  ] ", l.Ldate|l.Ltime|l.Lshortfile)
	Warn = l.New(file, "[ WARN  ] ", l.Ldate|l.Ltime|l.Lshortfile)
	Error = l.New(file, "[ ERROR ] ", l.Ldate|l.Ltime|l.Lshortfile)
}
