package echo

import (
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	at := time.Now()
	std := at.Format(time.RFC3339)
	var buf [26]byte
	n := FormatTimeRFC3339(buf[:], at)
	my := string(buf[:n])
	t.Log("my:", my)
	if std != my {
		t.Fail()
	}
}

func TestTrimSpace(t *testing.T) {
	t.Log(trimspace([]byte("")))
	t.Log(trimspace([]byte(" ")))
	t.Log(trimspace([]byte("a ")))
	t.Log(trimspace([]byte(" a")))
	t.Log(trimspace([]byte(" a\t")))
}
