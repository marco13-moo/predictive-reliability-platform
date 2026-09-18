package fixtures

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"time"
)

func Points() []contracts.MetricPoint {
	t := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	p := []contracts.MetricPoint{}
	for i := 0; i < 10; i++ {
		p = append(p, contracts.MetricPoint{Timestamp: t.Add(time.Duration(i) * time.Minute), Value: float64(i % 2)})
	}
	p = append(p, contracts.MetricPoint{Timestamp: t.Add(10 * time.Minute), Value: 20})
	return p
}
func Events() []contracts.Event {
	t := time.Date(2026, 1, 1, 0, 9, 0, 0, time.UTC)
	return []contracts.Event{{ID: "chg-1", Kind: contracts.KindChange, Service: "payments", Timestamp: t, Name: "deploy", Message: "v2"}}
}
