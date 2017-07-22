package echo

import (
	"fmt"
	"math"
	"runtime"

	"github.com/vizee/litebuf"
)

type FieldType uint

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
	TypeError
	TypeStringer
	TypeFormat
	TypeVar
	TypeEchoer
	TypeStack
	maxType
)

type Echoer interface {
	Echo(w *litebuf.Buffer)
}

type Field struct {
	Type FieldType
	Key  string
	Data interface{}
	U64  uint64
	Str  string
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
	return Field{Type: TypeString, Key: key, Str: val}
}

func Quote(key string, val string) Field {
	return Field{Type: TypeQuote, Key: key, Str: val}
}

func Errors(key string, val error) Field {
	return Field{Type: TypeError, Key: key, Data: val}
}

func Stringer(key string, val fmt.Stringer) Field {
	return Field{Type: TypeStringer, Key: key, Data: val}
}

func Var(key string, val interface{}) Field {
	return Field{Type: TypeVar, Key: key, Data: val}
}

func Echo(key string, val Echoer) Field {
	return Field{Type: TypeEchoer, Key: key, Data: val}
}

func Format(key string, format string, args ...interface{}) Field {
	return Field{Type: TypeFormat, Key: key, Str: format, Data: args}
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
	return Field{Type: TypeStack, Data: buf[:n]}
}
