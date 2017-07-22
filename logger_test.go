package echo

import (
	"errors"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func printStruct(t *testing.T, v interface{}) {
	rt := reflect.TypeOf(v)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t.Logf("%-20s %-10d %d", f.Name, f.Offset, f.Type.Size())
	}
}

func TestLoggerStruct(t *testing.T) {
	printStruct(t, Logger{})
}

func TestLog(t *testing.T) {
	l := Logger{}
	l.SetLevel(DebugLevel)
	fields := []Field{
		Field{Key: "nil"},
		Int("int", 1),
		Uint("uint", 2),
		Bool("bool", true),
		Hex("hex", 16),
		Float32("float32", 3.0),
		Float64("float64", 4.0),
		String("string", "string5"),
		Stringer("stringer", time.Now()),
		Errors("errors", errors.New("err")),
		Var("var", map[string]int{"a": 1, "b": 2}),
		Stack(false),
	}
	l.Debug("test", fields...)
}

func BenchmarkLog1Field(b *testing.B) {
	l := Logger{}
	l.SetLevel(DebugLevel)
	l.SetOutput(ioutil.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", String("key", "value"))
	}
}
