package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Level uint8

const (
	Debug Level = iota
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
	MinLevel       Level
	Folder         string
	FileNamePrefix string
	Period         time.Duration
	wg             sync.WaitGroup
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
	l.wg.Add(1)
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
	l.wg.Done()
}

func (l *Logger) Close() {
	close(l.ch)
	l.wg.Wait()

	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func (l *Logger) write(level Level, s string) {
	sl := strLevel[level]
	if l.WithFile {
		_, file, line, ok := runtime.Caller(3)
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
}

func (l *Logger) writeln(level Level, v ...interface{}) {
	if level >= l.MinLevel {
		l.write(level, fmt.Sprintln(v...))
	}
	if level == Fatal {
		l.Close()
		os.Exit(1)
	}
}

func (l *Logger) writef(level Level, format string, v ...interface{}) {
	if level >= l.MinLevel {
		l.write(level, fmt.Sprintf(format+"\n", v...))
	}
	if level == Fatal {
		l.Close()
		os.Exit(1)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	l.writeln(Debug, v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.writeln(Info, v...)
}

func (l *Logger) Warningln(v ...interface{}) {
	l.writeln(Warning, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.writeln(Error, v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.writeln(Fatal, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writef(Debug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.writef(Info, format, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.writef(Warning, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writef(Error, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.writef(Fatal, format, v...)
}

var Default Logger

func Debugln(v ...interface{}) {
	Default.writeln(Debug, v...)
}

func Infoln(v ...interface{}) {
	Default.writeln(Info, v...)
}

func Warningln(v ...interface{}) {
	Default.writeln(Warning, v...)
}

func Errorln(v ...interface{}) {
	Default.writeln(Error, v...)
}

func Fatalln(v ...interface{}) {
	Default.writeln(Fatal, v...)
}

func Debugf(format string, v ...interface{}) {
	Default.writef(Debug, format, v...)
}

func Infof(format string, v ...interface{}) {
	Default.writef(Info, format, v...)
}

func Warningf(format string, v ...interface{}) {
	Default.writef(Warning, format, v...)
}

func Errorf(format string, v ...interface{}) {
	Default.writef(Error, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Default.writef(Fatal, format, v...)
}
