package ingest

import (
	"fmt"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"sort"
	"strings"
	"time"
)

type Raw struct {
	ID, Kind, Source, Service, TenantID, Name, Message string
	SchemaVersion                                      int
	Timestamp                                          string
	Value                                              float64
	Labels                                             map[string]string
	Provenance                                         contracts.Provenance
}

func Normalize(r Raw) (contracts.Event, error) {
	id := strings.TrimSpace(r.ID)
	if id == "" || strings.TrimSpace(r.Service) == "" || strings.TrimSpace(r.TenantID) == "" {
		return contracts.Event{}, fmt.Errorf("id, service, and tenant id are required")
	}
	t, err := time.Parse(time.RFC3339Nano, r.Timestamp)
	if err != nil {
		return contracts.Event{}, fmt.Errorf("timestamp: %w", err)
	}
	k := contracts.EventKind(strings.ToLower(strings.TrimSpace(r.Kind)))
	if k != contracts.KindMetric && k != contracts.KindLog && k != contracts.KindChange {
		return contracts.Event{}, fmt.Errorf("unsupported kind %q", r.Kind)
	}
	v := r.SchemaVersion
	if v == 0 {
		v = 1
	}
	labels := map[string]string{}
	for key, val := range r.Labels {
		normalizedKey := strings.ToLower(strings.TrimSpace(key))
		if normalizedKey == "" {
			return contracts.Event{}, fmt.Errorf("label key is required")
		}
		if _, exists := labels[normalizedKey]; exists {
			return contracts.Event{}, fmt.Errorf("duplicate label key after normalization: %q", normalizedKey)
		}
		labels[normalizedKey] = strings.TrimSpace(val)
	}
	p := r.Provenance
	if strings.TrimSpace(p.Connector) == "" {
		p.Connector = strings.TrimSpace(r.Source)
		if p.Connector == "" {
			p.Connector = "unknown"
		}
	}
	if strings.TrimSpace(p.ParserVersion) == "" {
		p.ParserVersion = "normalize/v1"
	}
	e := contracts.Event{ID: id, SchemaVersion: v, TenantID: strings.TrimSpace(r.TenantID), Kind: k, Source: strings.TrimSpace(r.Source), Service: strings.TrimSpace(r.Service), Timestamp: t.UTC(), Name: strings.TrimSpace(r.Name), Unit: "count", Value: r.Value, Labels: labels, Message: strings.TrimSpace(r.Message), Provenance: p}
	return e, e.Validate()
}
func Batch(raw []Raw) ([]contracts.Event, []error) {
	out := make([]contracts.Event, 0, len(raw))
	errs := []error{}
	seen := map[string]bool{}
	for _, r := range raw {
		e, err := Normalize(r)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if seen[e.ID] {
			continue
		}
		seen[e.ID] = true
		out = append(out, e)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Timestamp.Before(out[j].Timestamp) })
	return out, errs
}
