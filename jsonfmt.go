package echo

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/vizee/litebuf"
)

type TimeStyle int

const (
	SimpleTime TimeStyle = iota
	TimeWithMs
	RFC3339Time
	UnixTimeStamp
	UnixTimeStampNano
)

var jsonTags = [...]string{
	FatalLevel: `"fatal"`,
	ErrorLevel: `"error"`,
	WarnLevel:  `"warn"`,
	InfoLevel:  `"info"`,
	DebugLevel: `"debug"`,
}

type JSONFormatter struct {
	TimeStyle     TimeStyle
	LevelTag      bool
	EscapeUnicode bool
}

func (f *JSONFormatter) Format(buf *litebuf.Buffer, at time.Time, level LogLevel, msg string, fields []Field) {
	buf.WriteString(`{"time":`)

	switch f.TimeStyle {
	case SimpleTime:
		buf.WriteByte('"')
		TimeFormat(buf.Reserve(19), at, false)
		buf.WriteByte('"')
	case TimeWithMs:
		buf.WriteByte('"')
		TimeFormat(buf.Reserve(23), at, true)
		buf.WriteByte('"')
	case RFC3339Time:
		buf.WriteByte('"')
		at.AppendFormat(buf.Reserve(25)[:0], time.RFC3339)
		buf.WriteByte('"')
	case UnixTimeStamp:
		buf.AppendInt(at.Unix(), 10)
	case UnixTimeStampNano:
		buf.AppendInt(at.UnixNano(), 10)
	default:
		panic(fmt.Sprintf("unknown time style: %d", f.TimeStyle))
	}

	buf.WriteString(`,"level":`)
	if f.LevelTag {
		buf.WriteString(jsonTags[level])
	} else {
		buf.AppendInt(int64(level), 10)
	}

	if msg != "" {
		buf.WriteString(`,"msg":`)
		buf.WriteQuote(msg, f.EscapeUnicode)
	}

	if len(fields) > 0 {
		buf.WriteString(`,"fields":{`)

		for i := range fields {
			if i > 0 {
				buf.WriteByte(',')
			}
			field := &fields[i]
			buf.WriteQuote(field.Key, f.EscapeUnicode)
			buf.WriteByte(':')

			switch field.Type {
			case TypeNil:
				buf.WriteString("null")
			case TypeInt, TypeInt32, TypeInt64:
				buf.AppendInt(int64(field.U64), 10)
			case TypeUint, TypeUint32, TypeUint64:
				buf.AppendUint(field.U64, 10)
			case TypeHex:
				buf.WriteString(`"0x`)
				buf.AppendUint(field.U64, 16)
				buf.WriteByte('"')
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
			case TypeString, TypeQuote:
				buf.WriteQuote(field.Str, f.EscapeUnicode)
			case TypeError:
				if err, ok := field.Data.(error); ok {
					buf.WriteQuote(err.Error(), f.EscapeUnicode)
				} else {
					buf.WriteString(`""`)
				}
			case TypeStringer:
				buf.WriteQuote(field.Data.(fmt.Stringer).String(), f.EscapeUnicode)
			case TypeFormat:
				tmpbuf := bufpool.Get().(*litebuf.Buffer)
				tmpbuf.Reset()
				fmt.Fprintf(buf, field.Str, field.Data.([]interface{}))
				buf.WriteQuote(tmpbuf.String(), false)
				bufpool.Put(tmpbuf)
			case TypeEchoer:
				field.Data.(Echoer).Echo(buf)
			case TypeVar:
				json.NewEncoder(buf).Encode(field.Data)
			default:
				buf.WriteString(fmt.Sprintf(`"skipped-type-%d"`, field.Type))
			}
		}

		buf.WriteByte('}')
	}

	buf.WriteByte('}')
}
