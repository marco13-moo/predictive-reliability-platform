package replay

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/detector"
	"sort"
	"time"
)

type Metrics struct {
	TruePositives, FalsePositives, FalseNegatives   int
	Precision, Recall, F1, FalsePositiveRate        float64
	DetectionLatency                                time.Duration
	RemediationSafetyRate                           float64
	AlertVolume, LabeledIncidents, UnknownIntervals int
	FormulaVersion                                  string
}

func Evaluate(d detector.Detector, points []contracts.MetricPoint, incidents []time.Time) Metrics {
	got := d.Detect("service", "metric", points)
	used := map[int]bool{}
	tp := 0
	var latency time.Duration
	for _, a := range got {
		best := -1
		for i, t := range incidents {
			if !used[i] && abs(a.Timestamp.Sub(t)) <= 5*time.Minute {
				best = i
				break
			}
		}
		if best >= 0 {
			used[best] = true
			tp++
			l := abs(a.Timestamp.Sub(incidents[best]))
			latency += l
		}
	}
	fp := len(got) - tp
	fn := len(incidents) - tp
	m := Metrics{TruePositives: tp, FalsePositives: fp, FalseNegatives: fn, AlertVolume: len(got), LabeledIncidents: len(incidents), FormulaVersion: "replay/v1"}
	if tp+fp > 0 {
		m.Precision = float64(tp) / float64(tp+fp)
		m.FalsePositiveRate = float64(fp) / float64(tp+fp)
	}
	if tp+fn > 0 {
		m.Recall = float64(tp) / float64(tp+fn)
	}
	if m.Precision+m.Recall > 0 {
		m.F1 = 2 * m.Precision * m.Recall / (m.Precision + m.Recall)
	}
	if tp > 0 {
		m.DetectionLatency = latency / time.Duration(tp)
	}
	return m
}
func RemediationSafety(approved, executed []bool) float64 {
	if len(executed) == 0 {
		return 1
	}
	safe := 0
	for i, v := range executed {
		if v && i < len(approved) && approved[i] {
			safe++
		}
	}
	return float64(safe) / float64(len(executed))
}
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
func Sort(points []contracts.MetricPoint) {
	sort.SliceStable(points, func(i, j int) bool {
		if points[i].Timestamp.Equal(points[j].Timestamp) {
			return points[i].ID < points[j].ID
		}
		return points[i].Timestamp.Before(points[j].Timestamp)
	})
}
