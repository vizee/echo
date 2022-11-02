package echo

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	l := NewLogger(os.Stdout, &PlainFormatter{})
	l.SetLevel(DebugLevel)
	fields := []Field{
		{Key: "nil"},
		Int("int", 1),
		Uint("uint", 2),
		Bool("bool", true),
		Hex("hex", 16),
		Float32("float32", 3.0),
		Float64("float64", 4.0),
		String("string", "string5"),
		Stringers("stringer", time.Now()),
		Errval("errors", io.EOF),
		Var("var", map[string]int{"a": 1, "b": 2}),
		Stack(false),
	}
	l.Debug("test", fields...)
}

func BenchmarkLog1Field(b *testing.B) {
	l := NewLogger(io.Discard, &PlainFormatter{})
	l.SetLevel(DebugLevel)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", String("key", "value"))
	}
}
