package echo

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/vizee/litebuf"
)

type rtype struct {
	typ    string
	fields []string
}

var rtcache sync.Map

func makertype(rt reflect.Type) *rtype {
	n := rt.NumField()
	t := &rtype{
		typ:    rt.String(),
		fields: make([]string, n),
	}
	for i := 0; i < n; i++ {
		rf := rt.Field(i)
		t.fields[i] = rf.Name
		if rf.Type.Kind() == reflect.Struct {
			if _, ok := rtcache.Load(rf.Type); !ok {
				makertype(rf.Type)
			}
		}
	}
	rtcache.Store(rt, t)
	return t
}

func getrtype(rt reflect.Type) *rtype {
	t, ok := rtcache.Load(rt)
	if ok {
		return t.(*rtype)
	}
	return makertype(rt)
}

func dumpStruct(buf *litebuf.Buffer, rv reflect.Value) {
	rt := getrtype(rv.Type())
	buf.WriteByte('{')
	for i := 0; i < len(rt.fields); i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(rt.fields[i])
		buf.WriteByte(':')
		field := rv.Field(i)
		if field.Kind() == reflect.Interface && field.CanInterface() {
			echoVar(buf, field.Interface(), false)
		} else {
			dumpValue(buf, field)
		}
	}
	buf.WriteByte('}')
}

func dumpMap(buf *litebuf.Buffer, rv reflect.Value) {
	buf.WriteString(`map[`)
	keys := rv.MapKeys()
	for i := range keys {
		if i > 0 {
			buf.WriteByte(' ')
		}
		dumpValue(buf, keys[i])
		buf.WriteByte(':')
		t := rv.MapIndex(keys[i])
		if t.Kind() == reflect.Interface {
			echoVar(buf, t.Interface(), false)
		} else {
			dumpValue(buf, t)
		}
	}
	buf.WriteByte(']')
}

func dumpArray(buf *litebuf.Buffer, rv reflect.Value) {
	buf.WriteByte('[')
	n := rv.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		t := rv.Index(i)
		if t.Kind() == reflect.Interface && t.CanInterface() {
			echoVar(buf, t.Interface(), false)
		} else {
			dumpValue(buf, t)
		}
	}
	buf.WriteByte(']')
}

func dumpValue(buf *litebuf.Buffer, rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteUint(rv.Uint(), 10)
	case reflect.Ptr, reflect.UnsafePointer, reflect.Chan, reflect.Func:
		p := rv.Pointer()
		if p == 0 {
			buf.WriteString("<nil>")
		} else {
			buf.WriteString("0x")
			buf.WriteUint(uint64(p), 16)
		}
	case reflect.Bool:
		if rv.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case reflect.Float32:
		buf.WriteFloat(rv.Float(), 'g', -1, 32)
	case reflect.Float64:
		buf.WriteFloat(rv.Float(), 'g', -1, 64)
	case reflect.String:
		buf.WriteString(rv.String())
	case reflect.Interface:
		if rv.IsNil() {
			buf.WriteString("<nil>")
		} else {
			p := rv.InterfaceData()
			buf.WriteString("{0x")
			buf.WriteUint(uint64(p[0]), 16)
			buf.WriteString(",0x")
			buf.WriteUint(uint64(p[1]), 16)
			buf.WriteByte('}')
		}
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		buf.WriteFloat(real(c), 'g', -1, 64)
		buf.WriteByte('+')
		buf.WriteFloat(imag(c), 'g', -1, 64)
		buf.WriteByte('i')
	case reflect.Map:
		if rv.IsNil() {
			buf.WriteString("<nil>")
		} else {
			dumpMap(buf, rv)
		}
	case reflect.Struct:
		dumpStruct(buf, rv)
	case reflect.Array, reflect.Slice:
		dumpArray(buf, rv)
	default:
		buf.WriteString("<invalid kind>")
	}
}

func echoVar(buf *litebuf.Buffer, x any, ptr bool) {
	switch t := x.(type) {
	case nil:
		buf.WriteString("<nil>")
	case int:
		buf.WriteInt(int64(t), 10)
	case int8:
		buf.WriteInt(int64(t), 10)
	case int16:
		buf.WriteInt(int64(t), 10)
	case int32:
		buf.WriteInt(int64(t), 10)
	case int64:
		buf.WriteInt(t, 10)
	case uint:
		buf.WriteUint(uint64(t), 10)
	case uint8:
		buf.WriteUint(uint64(t), 10)
	case uint16:
		buf.WriteUint(uint64(t), 10)
	case uint32:
		buf.WriteUint(uint64(t), 10)
	case uint64:
		buf.WriteUint(t, 10)
	case uintptr:
		buf.WriteUint(uint64(t), 10)
	case float32:
		buf.WriteFloat(float64(t), 'g', -1, 32)
	case float64:
		buf.WriteFloat(t, 'g', -1, 64)
	case bool:
		if t {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case unsafe.Pointer:
		buf.WriteString("0x")
		buf.WriteUint(uint64(uintptr(t)), 16)
	case string:
		buf.WriteString(t)
	case reflect.Value:
		dumpValue(buf, t)
	default:
		rv := reflect.ValueOf(t)
		if ptr && rv.Kind() == reflect.Ptr {
			buf.WriteByte('&')
			dumpValue(buf, rv.Elem())
		} else {
			dumpValue(buf, rv)
		}
	}
}
