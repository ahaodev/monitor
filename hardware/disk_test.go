package hardware

import "testing"

func TestDisk(t *testing.T) {
	s := GetDiskInfoForDisplay()
	t.Log(s)
}
