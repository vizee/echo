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

func TimeFormat(b []byte, t time.Time, msp bool) {
	tb := b[:19]
	year, month, day := t.Date()
	tb[0], tb[1] = getdigit(byte(year / 100))
	tb[2], tb[3] = getdigit(byte(year % 100))
	tb[4] = '-'
	tb[5], tb[6] = getdigit(byte(month))
	tb[7] = '-'
	tb[8], tb[9] = getdigit(byte(day))
	tb[10] = '/'
	hour, min, sec := t.Clock()
	tb[11], tb[12] = getdigit(byte(hour))
	tb[13] = ':'
	tb[14], tb[15] = getdigit(byte(min))
	tb[16] = ':'
	tb[17], tb[18] = getdigit(byte(sec))
	if msp {
		msb := b[19:23]
		ms := t.Nanosecond() / 1000000
		msb[0] = '.'
		msb[1] = byte(ms/100) + '0'
		msb[2], msb[3] = getdigit(byte(ms % 100))
	}
}
