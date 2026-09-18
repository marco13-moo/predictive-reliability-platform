package correlate

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"testing"
	"time"
)

func TestWindow(t *testing.T) {
	at := time.Now()
	a := contracts.Anomaly{Service: "s", Timestamp: at}
	es := []contracts.Event{{Kind: contracts.KindChange, Service: "s", Timestamp: at.Add(2 * time.Minute)}, {Kind: contracts.KindChange, Service: "s", Timestamp: at.Add(10 * time.Minute)}}
	if len(Changes(a, es, 3*time.Minute).Changes) != 1 {
		t.Fatal("window mismatch")
	}
}
