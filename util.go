package echo

import (
	"sync"
	"time"

	"github.com/vizee/litebuf"
)

var bufpool = sync.Pool{}

func getBuf() *litebuf.Buffer {
	b := bufpool.Get()
	if b != nil {
		buf := b.(*litebuf.Buffer)
		buf.Reset()
		return buf
	}
	return &litebuf.Buffer{}
}

func getdigit(n byte) (byte, byte) {
	return n/10 + '0', n%10 + '0'
}

func FormatSimpleTime(buf []byte, t time.Time) {
	// 2006-01-02 15:04:05
	b := buf[:19]
	year, month, day := t.Date()
	b[0], b[1] = getdigit(byte(year / 100))
	b[2], b[3] = getdigit(byte(year % 100))
	b[4] = '-'
	b[5], b[6] = getdigit(byte(month))
	b[7] = '-'
	b[8], b[9] = getdigit(byte(day))
	b[10] = ' '
	hour, min, sec := t.Clock()
	b[11], b[12] = getdigit(byte(hour))
	b[13] = ':'
	b[14], b[15] = getdigit(byte(min))
	b[16] = ':'
	b[17], b[18] = getdigit(byte(sec))
}

func FormatTimeRFC3339(buf []byte, t time.Time) int {
	// 2006-01-02T15:04:05Z07:00
	b := buf[:26]
	year, month, day := t.Date()
	b[0], b[1] = getdigit(byte(year / 100))
	b[2], b[3] = getdigit(byte(year % 100))
	b[4] = '-'
	b[5], b[6] = getdigit(byte(month))
	b[7] = '-'
	b[8], b[9] = getdigit(byte(day))
	b[10] = 'T'
	hour, min, sec := t.Clock()
	b[11], b[12] = getdigit(byte(hour))
	b[13] = ':'
	b[14], b[15] = getdigit(byte(min))
	b[16] = ':'
	b[17], b[18] = getdigit(byte(sec))
	_, offset := t.Zone()
	if offset == 0 {
		b[19] = 'Z'
		return 20
	}
	if offset > 0 {
		b[19] = '+'
	} else {
		b[19] = '-'
		offset = -offset
	}
	offset /= 60
	b[20], b[21] = getdigit(byte(offset / 60))
	b[22] = ':'
	b[23], b[24] = getdigit(byte(offset % 60))
	return 25
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
