package echo

import (
	"io"
	"os"
)

var DefaultFormatter Formatter = (*PlainFormatter)(nil)

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

func D(msg string, fields ...Field) {
	export.Debug(msg, fields...)
}

func I(msg string, fields ...Field) {
	export.Info(msg, fields...)
}

func W(msg string, fields ...Field) {
	export.Warn(msg, fields...)
}

func E(msg string, fields ...Field) {
	export.Error(msg, fields...)
}

func F(msg string, fields ...Field) {
	export.Fatal(msg, fields...)
}
