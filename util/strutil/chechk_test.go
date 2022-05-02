package strutil

import "testing"

func TestIsNumeric(t *testing.T) {
	var i string = "1"
	numeric := IsNumeric(i)
	t.Log(numeric)

}
