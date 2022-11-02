package echo

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"unsafe"

	"github.com/vizee/litebuf"
)

type SubAnonymousType struct {
	SAT  string
	Bool bool
}

type AnonymousType struct {
	SubAnonymousType
	AT string
}

type SubdataType struct {
	Name string
	Int  int
}

type stringer struct {
	target string
	n      int
}

func (w *stringer) String() string {
	return w.target
}

type DataType struct {
	Int       int
	Uint      uint
	Uintptr   uintptr
	Ptr       unsafe.Pointer
	Pint      *int
	Bool      bool
	Byte      byte
	Rune      rune
	Uint32    uint32
	String    string
	Float32   float32
	Float64   float64
	Complex64 complex64
	Iface     any
	NilIface  any
	Slice     []any
	Array     [8]any
	MapSI     map[string]any
	MapIS     map[int]string
	MapNil    map[int]string
	Stringer0 *stringer
	Stringer1 Stringer
	inner     *DataType
	Subdata   SubdataType
	AnonymousType
	Func func()
	Chan chan struct{}
}

var (
	n    = 126
	data = DataType{
		Int:       -1,
		Uint:      2,
		Uintptr:   128,
		Ptr:       unsafe.Pointer(&n),
		Pint:      &n,
		Bool:      true,
		Byte:      ' ',
		Rune:      '哟',
		String:    "你好\t世界",
		Float32:   math.Pi,
		Float64:   math.Pi,
		Complex64: complex(1, 2),
		Iface: &SubdataType{
			Name: "Iface",
			Int:  len("Iface"),
		},
		Slice: []any{"slice", 0},
		Array: [8]any{"array", false},
		MapSI: map[string]any{
			"int":    1,
			"bool":   false,
			"string": "zzz",
			"nil":    nil,
			"map":    map[string]string{"a": "b"},
		},
		MapIS: map[int]string{
			1: "one",
			2: "two",
			3: "three",
		},
		MapNil: nil,
		Stringer0: &stringer{
			n:      2,
			target: "i am stringer",
		},
		Stringer1: &stringer{
			n:      1,
			target: "/path/to/file",
		},
		inner: new(DataType),
		Subdata: SubdataType{
			Name: "sub",
			Int:  len("sub"),
		},
		AnonymousType: AnonymousType{
			AT: "anonymous",
		},
		Func: func() {},
		Chan: make(chan struct{}),
	}
)

func Test_dumpValue(t *testing.T) {
	buf := litebuf.Buffer{}
	t.Logf("fmt: %+v", &data)
	dumpValue(&buf, reflect.ValueOf(&data))
	t.Log("echo:", buf.String())
}

func BenchmarkPrintfPlusV(b *testing.B) {
	buf := litebuf.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		fmt.Fprintf(&buf, "%+v", &data)
	}
}

func BenchmarkEchoVar(b *testing.B) {
	buf := litebuf.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		dumpValue(&buf, reflect.ValueOf(&data))
	}
}
