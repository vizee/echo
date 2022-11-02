package echo

import "testing"

func TestExport(t *testing.T) {
	D("debug")
	I("info")
	W("warn")
	E("error")
}
