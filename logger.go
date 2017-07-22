package echo

import (
	"io"
	"os"
	"time"

	"github.com/vizee/litebuf"
)

type LogLevel uint

const (
	FatalLevel LogLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	maxLevel
)

type Formatter interface {
	Format(buf *litebuf.Buffer, t time.Time, level LogLevel, msg string, fields []Field)
}

var DefaultFormatter = &PlainFormatter{}

type Logger struct {
	level LogLevel
	fmter Formatter
	w     io.Writer
}

func (l *Logger) Level() LogLevel {
	return l.level
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetFormmatter(f Formatter) {
	l.fmter = f
}

func (l *Logger) SetOutput(w io.Writer) {
	l.w = w
}

func (l *Logger) log(level LogLevel, msg string, fields []Field) {
	buf := bufpool.Get().(*litebuf.Buffer)
	buf.Reset()

	f := l.fmter
	if f == nil {
		f = DefaultFormatter
	}
	f.Format(buf, time.Now(), level, msg, fields)
	buf.WriteByte('\n')

	w := l.w
	if w == nil {
		w = os.Stdout
	}
	w.Write(buf.Bytes())

	bufpool.Put(buf)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	if l.level >= DebugLevel {
		l.log(DebugLevel, msg, fields)
	}
}

func (l *Logger) Info(msg string, fields ...Field) {
	if l.level >= InfoLevel {
		l.log(InfoLevel, msg, fields)
	}
}

func (l *Logger) Warn(msg string, fields ...Field) {
	if l.level >= WarnLevel {
		l.log(WarnLevel, msg, fields)
	}
}

func (l *Logger) Error(msg string, fields ...Field) {
	if l.level >= ErrorLevel {
		l.log(ErrorLevel, msg, fields)
	}
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields)
	os.Exit(1)
}
