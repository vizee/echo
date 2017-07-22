package echo

import (
	"sync"
	"time"

	"github.com/vizee/litebuf"
)

type BytesEchoer []byte

func (be BytesEchoer) Echo(buf *litebuf.Buffer) {
	buf.Write([]byte(be))
}

var bufpool = sync.Pool{
	New: func() interface{} {
		return &litebuf.Buffer{}
	},
}

func getdigit(n byte) (byte, byte) {
	return n/10 + '0', n%10 + '0'
}

func TimeFormat(buf []byte, t time.Time) {
	buf = buf[:19]
	year, month, day := t.Date()
	buf[0], buf[1] = getdigit(byte(year / 100))
	buf[2], buf[3] = getdigit(byte(year % 100))
	buf[4] = '-'
	buf[5], buf[6] = getdigit(byte(month))
	buf[7] = '-'
	buf[8], buf[9] = getdigit(byte(day))
	buf[10] = '/'
	hour, min, sec := t.Clock()
	buf[11], buf[12] = getdigit(byte(hour))
	buf[13] = ':'
	buf[14], buf[15] = getdigit(byte(min))
	buf[16] = ':'
	buf[17], buf[18] = getdigit(byte(sec))
}
