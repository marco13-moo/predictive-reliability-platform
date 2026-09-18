package ingest

import "testing"

func TestNormalize(t *testing.T) {
	e, err := Normalize(Raw{ID: "1", Kind: "METRIC", Service: "api", TenantID: "t1", Timestamp: "2026-01-01T00:00:00Z"})
	if err != nil || e.Kind != "metric" || e.Timestamp.Location().String() != "UTC" {
		t.Fatalf("%+v %v", e, err)
	}
}
func TestBatchOrderingAndIdempotency(t *testing.T) {
	raw := []Raw{{ID: "2", Kind: "metric", Service: "a", TenantID: "t", Timestamp: "2026-01-01T00:01:00Z"}, {ID: "1", Kind: "metric", Service: "a", TenantID: "t", Timestamp: "2026-01-01T00:00:00Z"}, {ID: "1", Kind: "metric", Service: "a", TenantID: "t", Timestamp: "2026-01-01T00:00:00Z"}}
	got, errs := Batch(raw)
	if len(errs) != 0 || len(got) != 2 || got[0].ID != "1" {
		t.Fatalf("got=%v errs=%v", got, errs)
	}
}
