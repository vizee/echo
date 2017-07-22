package echo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/vizee/litebuf"
)

func TestJSONFormat(t *testing.T) {
	var stub interface{}
	buf := litebuf.Buffer{}
	f := JSONFormatter{}

	f.Format(&buf, time.Now(), DebugLevel, "Debug Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	err := json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.Format(&buf, time.Now(), InfoLevel, "Info Hello", []Field{String("who", "世\t界")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.EscapeUnicode = true
	f.Format(&buf, time.Now(), InfoLevel, "Info Hello", []Field{String("who", "世界")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.LevelTag = true
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.TimeStyle = RFC3339Time
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.TimeStyle = UnixTimeStamp
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	f.TimeStyle = UnixTimeStampNano
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	err = json.Unmarshal(buf.Bytes(), &stub)
	if err != nil {
		t.Error(err)
	}
}
