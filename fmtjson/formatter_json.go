package echo

import (
	"encoding/json"
	"math"
	"time"

	"github.com/vizee/echo"
	"github.com/vizee/litebuf"
)

type TimeStyle int

const (
	SimpleTime TimeStyle = iota
	RFC3339
	RFC3339Nano
	UnixTimestamp
	UnixTimestampNano
)

var jsonTags = [...]string{
	echo.FatalLevel: `"fatal"`,
	echo.ErrorLevel: `"error"`,
	echo.WarnLevel:  `"warn"`,
	echo.InfoLevel:  `"info"`,
	echo.DebugLevel: `"debug"`,
}

type JSONFormatter struct {
	TimeStyle     TimeStyle
	LevelTag      bool
	EscapeUnicode bool
}

func (f *JSONFormatter) Format(buf *litebuf.Buffer, at time.Time, level echo.LogLevel, msg string, fields []echo.Field) {
	buf.WriteString(`{"time":`)

	switch f.TimeStyle {
	default:
		buf.WriteByte('"')
		echo.FormatSimpleTime(buf.Reserve(19), at)
		buf.WriteByte('"')
	case RFC3339:
		buf.WriteByte('"')
		n := echo.FormatTimeRFC3339(buf.Reserve(25), at)
		buf.Trim(25 - n)
		buf.WriteByte('"')
	case RFC3339Nano:
		buf.WriteByte('"')
		n := len(at.AppendFormat(buf.Reserve(35)[:0], time.RFC3339Nano))
		buf.Trim(35 - n)
		buf.WriteByte('"')
	case UnixTimestamp:
		buf.WriteInt(at.Unix(), 10)
	case UnixTimestampNano:
		buf.WriteInt(at.UnixNano(), 10)
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
			case echo.TypeNil:
				buf.WriteString("null")
			case echo.TypeInt, echo.TypeInt32, echo.TypeInt64:
				buf.WriteInt(int64(field.U64), 10)
			case echo.TypeUint, echo.TypeUint32, echo.TypeUint64:
				buf.WriteUint(field.U64, 10)
			case echo.TypeHex:
				buf.WriteString(`"0x`)
				buf.WriteUint(field.U64, 16)
				buf.WriteByte('"')
			case echo.TypeBool:
				if field.U64 == 1 {
					buf.WriteString("true")
				} else {
					buf.WriteString("false")
				}
			case echo.TypeFloat32:
				buf.WriteFloat(float64(math.Float32frombits(uint32(field.U64))), 'f', -1, 32)
			case echo.TypeFloat64:
				buf.WriteFloat(math.Float64frombits(field.U64), 'f', -1, 64)
			case echo.TypeString, echo.TypeQuote:
				buf.WriteQuote(field.ToString(), f.EscapeUnicode)
			case echo.TypeError:
				e := field.ToError()
				if e != nil {
					buf.WriteQuote(e.Error(), f.EscapeUnicode)
				} else {
					buf.WriteString("null")
				}
			case echo.TypeStringer:
				buf.WriteQuote(field.ToStringer().String(), f.EscapeUnicode)
			case echo.TypeVar:
				data, err := json.Marshal(field.ToAny())
				if err != nil {
					buf.WriteString("{}")
					break
				}
				buf.Write(trimspace(data))
			case echo.TypeEcho:
				field.ToEcho().Echo(buf)
			default:
				buf.WriteString("{}")
			}
		}

		buf.WriteByte('}')
	}

	buf.WriteByte('}')
}

func isspace(c byte) bool {
	switch c {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

func trimspace(b []byte) []byte {
	l := 0
	for l < len(b) && isspace(b[l]) {
		l++
	}
	r := len(b) - 1
	for r > l && isspace(b[r]) {
		r--
	}
	return b[l : r+1]
}
