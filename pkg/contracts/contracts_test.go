package contracts

import (
	"testing"
	"time"
)

func TestEventValidationRejectsUnsupportedAndUnprovenanced(t *testing.T) {
	e := Event{ID: "1", SchemaVersion: 2, TenantID: "t", Kind: KindMetric, Timestamp: time.Now().UTC(), Unit: "count"}
	if err := e.Validate(); err == nil {
		t.Fatal("expected unsupported schema rejection")
	}
	e.SchemaVersion = 1
	if err := e.Validate(); err == nil {
		t.Fatal("expected provenance rejection")
	}
}
