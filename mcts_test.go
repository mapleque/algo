package algo

import "testing"

func TestMCTS(t *testing.T) {
	root := NewTree("", BoardSizeMini)
	root.expand()
	if len(root.children) != 25 {
		t.Error("first avialable steps should be 25, but:", len(root.children))
	}
}
