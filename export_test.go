package echo

import "testing"

func TestExport(t *testing.T) {
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
}
