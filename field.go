package echo

import (
	"math"
	"runtime"
	"unsafe"

	"github.com/vizee/litebuf"
)

type FieldType uint8

const (
	TypeNil FieldType = iota
	TypeBool
	TypeInt
	TypeInt32
	TypeInt64
	TypeUint
	TypeUint32
	TypeUint64
	TypeHex
	TypeFloat32
	TypeFloat64
	TypeString
	TypeQuote
	TypeBytes
	TypeError
	TypeStringer
	TypeVar
	TypeEcho
	TypeStack
)

type Value interface {
	Echo(w *litebuf.Buffer)
}

type Stringer interface {
	String() string
}

type Field struct {
	Key  string
	Ptr1 unsafe.Pointer
	Ptr2 unsafe.Pointer
	U64  uint64
	Type FieldType
}

func (f *Field) ToString() string {
	return raw2str(f.Ptr1, f.U64)
}

func (f *Field) ToBytes() []byte {
	return raw2bytes(f.Ptr1, f.U64)
}

func (f *Field) ToAny() any {
	return raw2any(f.Ptr1, f.Ptr2)
}

func (f *Field) ToError() error {
	return raw2iface[error](f.Ptr1, f.Ptr2)
}

func (f *Field) ToStringer() Stringer {
	return raw2iface[Stringer](f.Ptr1, f.Ptr2)
}

func (f *Field) ToEcho() Value {
	return raw2iface[Value](f.Ptr1, f.Ptr2)
}

func Bool(key string, val bool) Field {
	v := uint64(0)
	if val {
		v = 1
	}
	return Field{Type: TypeBool, Key: key, U64: v}
}

func Int(key string, val int) Field {
	return Field{Type: TypeInt, Key: key, U64: uint64(val)}
}

func Int32(key string, val int32) Field {
	return Field{Type: TypeInt32, Key: key, U64: uint64(val)}
}

func Int64(key string, val int64) Field {
	return Field{Type: TypeInt64, Key: key, U64: uint64(val)}
}

func Uint(key string, val uint) Field {
	return Field{Type: TypeUint, Key: key, U64: uint64(val)}
}

func Uint32(key string, val uint32) Field {
	return Field{Type: TypeUint32, Key: key, U64: uint64(val)}
}

func Uint64(key string, val uint64) Field {
	return Field{Type: TypeUint64, Key: key, U64: val}
}

func Hex(key string, val uintptr) Field {
	return Field{Type: TypeHex, Key: key, U64: uint64(val)}
}

func Float32(key string, val float32) Field {
	return Field{Type: TypeFloat32, Key: key, U64: uint64(math.Float32bits(val))}
}

func Float64(key string, val float64) Field {
	return Field{Type: TypeFloat64, Key: key, U64: math.Float64bits(val)}
}

func String(key string, val string) Field {
	p, n := str2raw(val)
	return Field{Type: TypeString, Key: key, Ptr1: p, U64: n}
}

func Quote(key string, val string) Field {
	p, n := str2raw(val)
	return Field{Type: TypeQuote, Key: key, Ptr1: p, U64: n}
}

func Bytes(key string, val []byte) Field {
	p, n := bytes2raw(val)
	return Field{Type: TypeBytes, Key: key, Ptr1: p, U64: n}
}

func Errval(key string, val error) Field {
	p1, p2 := iface2raw(val)
	return Field{Type: TypeError, Key: key, Ptr1: p1, Ptr2: p2}
}

func Errors(val error) Field {
	return Errval("error", val)
}

func Stringers(key string, val Stringer) Field {
	p1, p2 := iface2raw(val)
	return Field{Type: TypeStringer, Key: key, Ptr1: p1, Ptr2: p2}
}

func Var(key string, val any) Field {
	p1, p2 := any2raw(val)
	return Field{Type: TypeVar, Key: key, Ptr1: p1, Ptr2: p2}
}

func Echo(key string, val Value) Field {
	p1, p2 := iface2raw(val)
	return Field{Type: TypeEcho, Key: key, Ptr1: p1, Ptr2: p2}
}

func Stack(all bool) Field {
	n := 1 << 12
	if all {
		n <<= 8
	}
	var buf []byte
	for n <= 64<<20 {
		buf = make([]byte, n)
		n = runtime.Stack(buf, all)
		if n < len(buf) {
			break
		}
		n += n
	}
	p, m := bytes2raw(buf)
	return Field{Type: TypeStack, Ptr1: p, U64: m}
}
