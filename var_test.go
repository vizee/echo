package echo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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
	Iface     interface{}
	NilIface  interface{}
	Slice     []interface{}
	Array     [8]interface{}
	MapSI     map[string]interface{}
	MapIS     map[int]string
	MapNil    map[int]string
	Stringer  fmt.Stringer
	inner     *DataType
	Subdata   SubdataType
	Encoder   *json.Encoder
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
		Slice: []interface{}{"slice", 0},
		Array: [8]interface{}{"array", false},
		MapSI: map[string]interface{}{
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
		Stringer: &stringer{
			n:      1,
			target: "/path/to/file",
		},
		inner: new(DataType),
		Subdata: SubdataType{
			Name: "sub",
			Int:  len("sub"),
		},
		Encoder: &json.Encoder{},
		AnonymousType: AnonymousType{
			AT: "anonymous",
		},
		Func: func() {},
		Chan: make(chan struct{}),
	}
)

func TestEchoVar(t *testing.T) {
	buf := litebuf.Buffer{}
	echoVar(&buf, &data, true)
	t.Logf("%+v", &data)
	t.Log(buf.String())
}

func BenchmarkPrintfPlusV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(ioutil.Discard, "%+v", &data)
	}
}

func BenchmarkEchoVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := litebuf.Buffer{}
		echoVar(&buf, &data, true)
	}
}
