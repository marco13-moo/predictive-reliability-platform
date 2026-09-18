package audit

import "testing"

func TestTrail(t *testing.T) {
	var tr Trail
	_, err := tr.Append("a", "x", map[string]int{"n": 1})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tr.Append("a", "y", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !tr.Verify() {
		t.Fatal("invalid")
	}
}
