package correlate

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"math"
	"time"
)

func Changes(a contracts.Anomaly, events []contracts.Event, window time.Duration) contracts.Correlation {
	c := contracts.Correlation{Anomaly: a}
	for _, e := range events {
		if e.Kind == contracts.KindChange && e.Service == a.Service && math.Abs(e.Timestamp.Sub(a.Timestamp).Seconds()) <= window.Seconds() {
			c.Changes = append(c.Changes, e)
		}
	}
	if len(c.Changes) > 0 {
		c.Score = math.Min(1, float64(len(c.Changes))/3)
	}
	return c
}
