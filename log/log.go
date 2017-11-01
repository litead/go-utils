package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	Debug = iota
	Info
	Warning
	Error
	Fatal
)

var strLevel = []string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"FATAL",
}

type Logger struct {
	ToStdErr       bool
	ToFile         bool
	WithFile       bool
	NoTime         bool
	MinLevel       uint8
	Folder         string
	FileNamePrefix string
	Period         time.Duration
	ch             chan string
	file           *os.File
	fileExpireAt   time.Time
}

func (l *Logger) Start() error {
	l.ch = make(chan string, 256)

	l.Period = l.Period.Round(time.Minute)
	if l.Period < time.Minute {
		l.Period = 24 * time.Hour
	}

	now := time.Now()
	y, m, d := now.Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	l.fileExpireAt = day.Add(now.Sub(day) / l.Period * l.Period)

	go l.run()
	return nil
}

func (l *Logger) createNewFile() error {
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}

	path := l.FileNamePrefix + l.fileExpireAt.Format("20060102_1504.log")
	path = filepath.Join(l.Folder, path)
	if f, e := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); e != nil {
		return e
	} else {
		l.file = f
	}

	l.fileExpireAt = l.fileExpireAt.Add(l.Period)
	return nil
}

func (l *Logger) run() {
	for s, ok := <-l.ch; ok; s, ok = <-l.ch {
		now := time.Now()
		if !l.NoTime {
			s = now.Format("15:04:05.000\t") + s
		}
		if l.ToStdErr {
			os.Stderr.WriteString(s)
		}
		if !l.ToFile {
			continue
		}
		if !now.Before(l.fileExpireAt) {
			if e := l.createNewFile(); e != nil {
				continue
			}
		}
		l.file.WriteString(s)
	}
}

func (l *Logger) Close() {
	close(l.ch)

	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func (l *Logger) write(level uint8, s string) {
	sl := strLevel[level]
	if l.WithFile {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		} else {
			_, file = filepath.Split(file)
		}
		s = fmt.Sprintf("%s\t%s(%d)\t%s", sl, file, line, s)
	} else {
		s = sl + "\t" + s
	}

	l.ch <- s
	if level == Fatal {
		l.Close()
		os.Exit(1)
	}
}

func (l *Logger) Write(level uint8, v ...interface{}) {
	if level >= l.MinLevel {
		l.write(level, fmt.Sprintln(v...))
	}
}

func (l *Logger) Writef(level uint8, format string, v ...interface{}) {
	if level >= l.MinLevel {
		l.write(level, fmt.Sprintf(format+"\n", v...))
	}
}
