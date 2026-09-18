package detector

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/fixtures"
	"testing"
)

func TestDetect(t *testing.T) {
	a := Detector{ZThreshold: 2}.Detect("x", "m", fixtures.Points())
	if len(a) != 1 {
		t.Fatalf("got %d", len(a))
	}
}
