package echo

import (
	"math"
	"reflect"
	"time"

	"github.com/vizee/litebuf"
)

var plainTags = [...]string{
	FatalLevel: " |F|",
	ErrorLevel: " |E|",
	WarnLevel:  " |W|",
	InfoLevel:  " |I|",
	DebugLevel: " |D|",
}

type PlainFormatter struct{}

func (*PlainFormatter) Format(buf *litebuf.Buffer, at time.Time, level LogLevel, msg string, fields []Field) {
	FormatSimpleTime(buf.Reserve(19), at)

	buf.WriteString(plainTags[level])

	if msg != "" {
		buf.WriteByte(' ')
		buf.WriteString(msg)
	}

	if len(fields) > 0 {
		buf.WriteString(" {")

		for i := range fields {
			if i > 0 {
				buf.WriteByte(' ')
			}
			field := &fields[i]
			if field.Key != "" {
				buf.WriteString(field.Key)
				buf.WriteByte('=')
			}

			switch field.Type {
			case TypeNil:
				buf.WriteString("nil")
			case TypeInt, TypeInt32, TypeInt64:
				buf.WriteInt(int64(field.U64), 10)
			case TypeUint, TypeUint32, TypeUint64:
				buf.WriteUint(field.U64, 10)
			case TypeHex:
				buf.WriteString("0x")
				buf.WriteUint(field.U64, 16)
			case TypeBool:
				if field.U64 == 1 {
					buf.WriteString("true")
				} else {
					buf.WriteString("false")
				}
			case TypeFloat32:
				buf.WriteFloat(float64(math.Float32frombits(uint32(field.U64))), 'f', -1, 32)
			case TypeFloat64:
				buf.WriteFloat(math.Float64frombits(field.U64), 'f', -1, 64)
			case TypeString:
				buf.WriteString(field.ToString())
			case TypeQuote:
				buf.WriteQuote(field.ToString(), false)
			case TypeBytes:
				buf.Write(field.ToBytes())
			case TypeError:
				e := field.ToError()
				if e != nil {
					buf.WriteString(e.Error())
				} else {
					buf.WriteString("nil")
				}
			case TypeStringer:
				buf.WriteString(field.ToStringer().String())
			case TypeVar:
				dumpValue(buf, reflect.ValueOf(field.ToAny()))
			case TypeEcho:
				field.ToEcho().Echo(buf)
			case TypeStack:
				buf.WriteString("stack: <\n")
				buf.Write(field.ToBytes())
				buf.WriteString(">\n")
			default:
				buf.WriteString("<invalid type>")
			}
		}

		buf.WriteByte('}')
	}
}
