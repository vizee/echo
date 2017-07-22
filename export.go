package echo

import (
	"io"
	"os"
)

var export = Logger{
	level: InfoLevel,
	w:     os.Stdout,
	fmter: DefaultFormatter,
}

func Level() LogLevel {
	return export.Level()
}

func SetLevel(level LogLevel) {
	export.SetLevel(level)
}

func SetOutput(w io.Writer) {
	export.SetOutput(w)
}

func SetFormatter(f Formatter) {
	export.SetFormmatter(f)
}

func Debug(msg string, fields ...Field) {
	export.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	export.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	export.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	export.Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	export.Fatal(msg, fields...)
}
