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

type Logger[W io.Writer, F Formatter] struct {
	level LogLevel
	fmter F
	w     W
}

func NewLogger[W io.Writer, F Formatter](w W, fmter F) Logger[W, F] {
	return Logger[W, F]{
		level: InfoLevel,
		fmter: fmter,
		w:     w,
	}
}

func (l *Logger[W, F]) Level() LogLevel {
	return l.level
}

func (l *Logger[W, F]) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger[W, F]) SetFormmatter(f F) {
	l.fmter = f
}

func (l *Logger[W, F]) SetOutput(w W) {
	l.w = w
}

func (l *Logger[W, F]) log(level LogLevel, msg string, fields []Field) {
	buf := getBuf()

	l.fmter.Format(buf, time.Now(), level, msg, fields)
	buf.WriteByte('\n')
	l.w.Write(buf.Bytes())

	bufpool.Put(buf)
}

func (l *Logger[W, F]) Debug(msg string, fields ...Field) {
	if l.level >= DebugLevel {
		l.log(DebugLevel, msg, fields)
	}
}

func (l *Logger[W, F]) Info(msg string, fields ...Field) {
	if l.level >= InfoLevel {
		l.log(InfoLevel, msg, fields)
	}
}

func (l *Logger[W, F]) Warn(msg string, fields ...Field) {
	if l.level >= WarnLevel {
		l.log(WarnLevel, msg, fields)
	}
}

func (l *Logger[W, F]) Error(msg string, fields ...Field) {
	if l.level >= ErrorLevel {
		l.log(ErrorLevel, msg, fields)
	}
}

type syncer interface {
	Sync() error
}

func (l *Logger[W, F]) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields)
	if s, ok := any(l.w).(syncer); ok {
		s.Sync()
	}
	os.Exit(1)
}
