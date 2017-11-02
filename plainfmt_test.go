package echo

import (
	"testing"
	"time"

	"github.com/vizee/litebuf"
)

func TestPlainFormat(t *testing.T) {
	buf := litebuf.Buffer{}
	f := PlainFormatter{}
	f.Format(&buf, time.Now(), DebugLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	buf.Reset()
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "世\t界")})
	t.Log(buf.String())
}
