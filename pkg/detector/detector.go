package detector

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"math"
	"sort"
)

type Detector struct {
	ZThreshold float64
	MinSamples int
}

func (d Detector) Detect(service, metric string, points []contracts.MetricPoint) []contracts.Anomaly {
	if d.ZThreshold == 0 {
		d.ZThreshold = 3
	}
	if d.MinSamples == 0 {
		d.MinSamples = 5
	}
	if len(points) < d.MinSamples {
		return nil
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Timestamp.Before(points[j].Timestamp) })
	vals := make([]float64, len(points))
	for i, p := range points {
		vals[i] = p.Value
	}
	mean := 0.0
	for _, v := range vals {
		mean += v
	}
	mean /= float64(len(vals))
	sd := 0.0
	for _, v := range vals {
		sd += (v - mean) * (v - mean)
	}
	sd = math.Sqrt(sd / float64(len(vals)))
	if sd == 0 {
		return nil
	}
	out := []contracts.Anomaly{}
	for _, p := range points {
		z := math.Abs((p.Value - mean) / sd)
		if z >= d.ZThreshold {
			out = append(out, contracts.Anomaly{Timestamp: p.Timestamp, Service: service, Metric: metric, Score: z, Value: p.Value, Baseline: mean, Reason: "z-score exceeded threshold"})
		}
	}
	return out
}
