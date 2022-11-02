package echo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/vizee/echo"
	"github.com/vizee/litebuf"
)

func TestJSONFormat(t *testing.T) {
	var stub any
	buf := litebuf.Buffer{}
	f := JSONFormatter{}

	f.Format(&buf, time.Now(), echo.DebugLevel, "Debug Hello", []echo.Field{echo.String("who", "World")})
	t.Log(buf.String())
	err := json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.Format(&buf, time.Now(), echo.InfoLevel, "Info Hello", []echo.Field{echo.String("who", "世\t界")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.EscapeUnicode = true
	f.Format(&buf, time.Now(), echo.InfoLevel, "Info Hello", []echo.Field{echo.String("who", "世界")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.LevelTag = true
	f.Format(&buf, time.Now(), echo.InfoLevel, "Hello", []echo.Field{echo.String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.TimeStyle = RFC3339Nano
	f.Format(&buf, time.Now(), echo.InfoLevel, "Hello", []echo.Field{echo.String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.TimeStyle = UnixTimestamp
	f.Format(&buf, time.Now(), echo.InfoLevel, "Hello", []echo.Field{echo.String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.TimeStyle = UnixTimestampNano
	f.Format(&buf, time.Now(), echo.InfoLevel, "Hello", []echo.Field{echo.String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}
}

func Test_trimspace(t *testing.T) {
	t.Log(trimspace([]byte("")))
	t.Log(trimspace([]byte(" ")))
	t.Log(trimspace([]byte("a ")))
	t.Log(trimspace([]byte(" a")))
	t.Log(trimspace([]byte(" a\t")))
}
