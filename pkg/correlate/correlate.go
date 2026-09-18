package correlate

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"math"
	"sort"
	"time"
)

func Changes(a contracts.Anomaly, events []contracts.Event, window time.Duration) contracts.Correlation {
	c := contracts.Correlation{Anomaly: a}
	for _, e := range events {
		if e.Kind == contracts.KindChange && e.Service == a.Service && (a.TenantID == "" || e.TenantID == a.TenantID) && math.Abs(e.Timestamp.Sub(a.Timestamp).Seconds()) <= window.Seconds() {
			c.Changes = append(c.Changes, e)
		}
	}
	sort.SliceStable(c.Changes, func(i, j int) bool {
		di, dj := math.Abs(c.Changes[i].Timestamp.Sub(a.Timestamp).Seconds()), math.Abs(c.Changes[j].Timestamp.Sub(a.Timestamp).Seconds())
		if di == dj {
			return c.Changes[i].ID < c.Changes[j].ID
		}
		return di < dj
	})
	if len(c.Changes) > 0 {
		c.Score = math.Min(1, float64(len(c.Changes))/3)
		c.Factors = []string{"tenant-scoped identity match", "event-time window proximity"}
	}
	c.Window = window
	return c
}
