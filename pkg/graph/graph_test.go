package graph

import (
	"testing"
	"time"
)

func TestBuilderIsTenantScopedAndDeterministic(t *testing.T) {
	asOf := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b := DeterministicBuilder{}
	g, err := b.Build("t1", asOf, []Relationship{
		{ID: "b", TenantID: "t1", From: "api", To: "db", ObservedAt: asOf, Confidence: .8},
		{ID: "a", TenantID: "t2", From: "secret", To: "db", ObservedAt: asOf},
		{ID: "b", TenantID: "t1", From: "api", To: "db", ObservedAt: asOf},
	})
	if err != nil || len(g.Edges) != 1 || g.Edges[0].TenantID != "t1" || g.Complete {
		t.Fatalf("unexpected graph: %+v, %v", g, err)
	}
}
