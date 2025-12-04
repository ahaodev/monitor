package hardware

import "testing"

func TestMemory(t *testing.T) {
	s := GetMemInfo()
	t.Log(s)
}
