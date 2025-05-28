package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var projectRoot string

func Init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working dir: %v", err)
	}
	projectRoot = wd

	log.SetFlags(0)
	log.SetPrefix("[product-service] ")
	log.SetOutput(&relWriter{Out: os.Stdout})
}

type relWriter struct {
	Out io.Writer
}

func (w *relWriter) Write(p []byte) (n int, err error) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "unknown"
		line = 0
	}
	rel, err := filepath.Rel(projectRoot, file)
	if err != nil {
		rel = file
	}
	header := fmt.Sprintf(
		"%s %s:%d: ",
		time.Now().Format("2006/01/02 15:04:05"),
		rel,
		line,
	)
	if _, err = w.Out.Write([]byte(header)); err != nil {
		return 0, err
	}
	return w.Out.Write(p)
}
