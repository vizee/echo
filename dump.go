package echo

import (
	"reflect"
	"unsafe"

	"github.com/vizee/litebuf"
)

func dumpAny(buf *litebuf.Buffer, rv reflect.Value, expand int) {
	switch t := rv.Interface().(type) {
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
	case Stringer:
		buf.WriteString(t.String())
	case Value:
		t.Echo(buf)
	case reflect.Value:
		dumpValue2(buf, t, expand)
	default:
		p := rv.InterfaceData()
		buf.WriteString("{0x")
		buf.WriteUint(uint64(p[0]), 16)
		buf.WriteString(",0x")
		buf.WriteUint(uint64(p[1]), 16)
		buf.WriteByte('}')
	}
}

func dumpStruct(buf *litebuf.Buffer, rv reflect.Value, expand int) {
	rt := rv.Type()
	buf.WriteByte('{')
	n := rv.NumField()
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(rt.Field(i).Name)
		buf.WriteByte(':')
		dumpValue2(buf, rv.Field(i), expand)
	}
	buf.WriteByte('}')
}

func dumpMap(buf *litebuf.Buffer, rv reflect.Value, expand int) {
	buf.WriteString("map[")
	keys := rv.MapKeys()
	for i := range keys {
		if i > 0 {
			buf.WriteByte(' ')
		}
		dumpValue2(buf, keys[i], expand)
		buf.WriteByte(':')
		dumpValue2(buf, rv.MapIndex(keys[i]), expand)
	}
	buf.WriteByte(']')
}

func dumpArray(buf *litebuf.Buffer, rv reflect.Value, expand int) {
	buf.WriteByte('[')
	n := rv.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		dumpValue2(buf, rv.Index(i), expand)
	}
	buf.WriteByte(']')
}

func dumpValue2(buf *litebuf.Buffer, rv reflect.Value, expand int) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteUint(rv.Uint(), 10)
	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		p := rv.Pointer()
		if p == 0 {
			buf.WriteString("<nil>")
		} else {
			buf.WriteString("0x")
			buf.WriteUint(uint64(p), 16)
		}
	case reflect.Ptr:
		if expand > 0 {
			buf.WriteByte('&')
			dumpValue2(buf, rv.Elem(), expand-1)
		} else {
			p := rv.Pointer()
			if p == 0 {
				buf.WriteString("<nil>")
			} else {
				buf.WriteString("0x")
				buf.WriteUint(uint64(p), 16)
			}
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
			dumpAny(buf, rv, expand)
		}
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		buf.WriteByte('(')
		buf.WriteFloat(real(c), 'g', -1, 64)
		buf.WriteByte('+')
		buf.WriteFloat(imag(c), 'g', -1, 64)
		buf.WriteByte('i')
		buf.WriteByte(')')
	case reflect.Map:
		if rv.IsNil() {
			buf.WriteString("<nil>")
		} else {
			dumpMap(buf, rv, expand)
		}
	case reflect.Struct:
		dumpStruct(buf, rv, expand)
	case reflect.Array, reflect.Slice:
		dumpArray(buf, rv, expand)
	default:
		buf.WriteString("<invalid kind>")
	}
}

func dumpValue(buf *litebuf.Buffer, rv reflect.Value) {
	dumpValue2(buf, rv, 1)
}
