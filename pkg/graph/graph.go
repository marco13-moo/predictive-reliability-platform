package graph

import (
	"fmt"
	"sort"
	"time"
)

type Relationship struct {
	ID, TenantID, From, To, Relation, Source string
	ObservedAt                               time.Time
	EffectiveFrom, EffectiveTo               time.Time
	Confidence                               float64
}
type Node struct{ ID, TenantID, Kind string }
type Edge struct {
	Relationship
	MissingEndpoint bool
}
type Graph struct {
	TenantID  string
	AsOf      time.Time
	Nodes     []Node
	Edges     []Edge
	Complete  bool
	Freshness time.Duration
}

type Builder interface {
	Build(tenant string, asOf time.Time, events []Relationship) (Graph, error)
}
type DeterministicBuilder struct{}

func (DeterministicBuilder) Build(tenant string, asOf time.Time, events []Relationship) (Graph, error) {
	if tenant == "" || asOf.IsZero() {
		return Graph{}, fmt.Errorf("tenant and as-of time are required")
	}
	g := Graph{TenantID: tenant, AsOf: asOf.UTC(), Complete: true}
	seen := map[string]bool{}
	for _, r := range events {
		if r.TenantID != tenant || r.ID == "" || r.From == "" || r.To == "" {
			continue
		}
		if r.ObservedAt.After(asOf) || (!r.EffectiveFrom.IsZero() && asOf.Before(r.EffectiveFrom)) || (!r.EffectiveTo.IsZero() && asOf.After(r.EffectiveTo)) {
			continue
		}
		if seen[r.ID] {
			continue
		}
		seen[r.ID] = true
		g.Edges = append(g.Edges, Edge{Relationship: r})
		g.Nodes = append(g.Nodes, Node{ID: r.From, TenantID: tenant}, Node{ID: r.To, TenantID: tenant})
	}
	sort.Slice(g.Edges, func(i, j int) bool { return g.Edges[i].ID < g.Edges[j].ID })
	sort.Slice(g.Nodes, func(i, j int) bool { return g.Nodes[i].ID < g.Nodes[j].ID })
	unique := g.Nodes[:0]
	for _, n := range g.Nodes {
		if len(unique) == 0 || unique[len(unique)-1].ID != n.ID {
			unique = append(unique, n)
		}
	}
	g.Nodes = unique
	g.Complete = len(g.Edges) == len(events)
	return g, nil
}
