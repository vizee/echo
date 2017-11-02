package echo

import (
	"fmt"
	"math"
	"time"

	"github.com/vizee/litebuf"
)

var plainTags = [...]string{
	FatalLevel: " [FAT]",
	ErrorLevel: " [ERR]",
	WarnLevel:  " [WRN]",
	InfoLevel:  " [INF]",
	DebugLevel: " [DBG]",
}

type PlainFormatter struct{}

func (*PlainFormatter) Format(buf *litebuf.Buffer, at time.Time, level LogLevel, msg string, fields []Field) {
	TimeFormat(buf.Reserve(23), at, true)

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
				buf.AppendInt(int64(field.U64), 10)
			case TypeUint, TypeUint32, TypeUint64:
				buf.AppendUint(field.U64, 10)
			case TypeHex:
				buf.WriteString("0x")
				buf.AppendUint(field.U64, 16)
			case TypeBool:
				if field.U64 == 1 {
					buf.WriteString("true")
				} else {
					buf.WriteString("false")
				}
			case TypeFloat32:
				buf.AppendFloat(float64(math.Float32frombits(uint32(field.U64))), 'f', -1, 32)
			case TypeFloat64:
				buf.AppendFloat(math.Float64frombits(field.U64), 'f', -1, 64)
			case TypeString:
				buf.WriteString(field.Str)
			case TypeQuote:
				buf.WriteQuote(field.Str, false)
			case TypeError:
				if err, ok := field.Data.(error); ok {
					buf.WriteString(err.Error())
				} else {
					buf.WriteString("nil")
				}
			case TypeStringer:
				buf.WriteString(field.Data.(fmt.Stringer).String())
			case TypeFormat:
				fmt.Fprintf(buf, field.Str, field.Data.([]interface{})...)
			case TypeVar:
				echoVar(buf, field.Data, true)
			case TypeEchoer:
				field.Data.(Echoer).Echo(buf)
			case TypeStack:
				buf.WriteByte('\n')
				buf.Write(field.Data.([]byte))
			default:
				buf.WriteString(fmt.Sprintf("skipped-type(%d)", field.Type))
			}
		}

		buf.WriteByte('}')
	}
}
