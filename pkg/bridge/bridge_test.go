package bridge

import "testing"

func TestBridge(t *testing.T) {
	w := Window{}
	t.Fatal(getPrefix(w))
}
