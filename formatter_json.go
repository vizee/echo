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
	RFC3339
	RFC3339Nano
	UnixTimestamp
	UnixNanoTimestamp
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
		FormatSimpleTime(buf.Reserve(19), at)
		buf.WriteByte('"')
	case RFC3339:
		buf.WriteByte('"')
		n := FormatTimeRFC3339(buf.Reserve(25), at)
		buf.Trim(25 - n)
		buf.WriteByte('"')
	case RFC3339Nano:
		buf.WriteByte('"')
		n := len(at.AppendFormat(buf.Reserve(35)[:0], time.RFC3339Nano))
		buf.Trim(35 - n)
		buf.WriteByte('"')
	case UnixTimestamp:
		buf.WriteInt(at.Unix(), 10)
	case UnixNanoTimestamp:
		buf.WriteInt(at.UnixNano(), 10)
	default:
		panic(fmt.Sprintf("unknown time style: %d", f.TimeStyle))
	}

	buf.WriteString(`,"level":`)
	if f.LevelTag {
		buf.WriteString(jsonTags[level])
	} else {
		buf.WriteInt(int64(level), 10)
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
				buf.WriteInt(int64(field.U64), 10)
			case TypeUint, TypeUint32, TypeUint64:
				buf.WriteUint(field.U64, 10)
			case TypeHex:
				buf.WriteString(`"0x`)
				buf.WriteUint(field.U64, 16)
				buf.WriteByte('"')
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
			case TypeString, TypeQuote:
				buf.WriteQuote(field.str(), f.EscapeUnicode)
			case TypeError:
				e := field.error()
				if e != nil {
					buf.WriteQuote(e.Error(), f.EscapeUnicode)
				} else {
					buf.WriteString("null")
				}
			case TypeStringer:
				buf.WriteQuote(field.stringer().String(), f.EscapeUnicode)
			case TypeFormat:
				tmpbuf := getBuf()
				fmtargs := field.fmtargs()
				fmt.Fprintf(tmpbuf, fmtargs.f, fmtargs.args...)
				buf.WriteQuote(tmpbuf.UnsafeString(), false)
				bufpool.Put(tmpbuf)
			case TypeVar:
				data, err := json.Marshal(field.Ptr1)
				if err != nil {
					buf.WriteString("nil")
					break
				}
				buf.Write(trimspace(data))
			case TypeEcho:
				field.echo().Echo(buf)
			default:
				buf.WriteString(fmt.Sprintf(`"skipped-type-%d"`, field.Type))
			}
		}

		buf.WriteByte('}')
	}

	buf.WriteByte('}')
}
