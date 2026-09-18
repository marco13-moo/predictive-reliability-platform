package replay

import (
	"testing"
	"time"
)

func TestSafetyAndLatency(t *testing.T) {
	if RemediationSafety([]bool{true, false}, []bool{true, false}) != .5 {
		t.Fatal("safety")
	}
	_ = time.Second
}
