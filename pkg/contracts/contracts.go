package contracts

import (
	"fmt"
	"strings"
	"time"
)

type EventKind string

const (
	KindMetric EventKind = "metric"
	KindLog    EventKind = "log"
	KindChange EventKind = "change"
)

type Provenance struct{ Connector, Region, ReceivedBy string }
type Event struct {
	ID            string
	SchemaVersion int
	TenantID      string
	Kind          EventKind
	Source        string
	Service       string
	Timestamp     time.Time
	Name          string
	Value         float64
	Labels        map[string]string
	Message       string
	Provenance    Provenance
}

func (e Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return fmt.Errorf("event id is required")
	}
	if e.SchemaVersion < 1 {
		return fmt.Errorf("schema version must be positive")
	}
	if strings.TrimSpace(e.TenantID) == "" {
		return fmt.Errorf("tenant scope is required")
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	if e.Kind != KindMetric && e.Kind != KindLog && e.Kind != KindChange {
		return fmt.Errorf("unsupported event kind %q", e.Kind)
	}
	return nil
}

type MetricPoint struct {
	ID        string
	Timestamp time.Time
	Value     float64
}
type SLO struct {
	Name, Service string
	Target        float64
	Window        time.Duration
}
type Budget struct {
	SLO       SLO
	Total     time.Duration
	Consumed  time.Duration
	Remaining time.Duration
	BurnRate  float64
	Exhausted bool
}
type Anomaly struct {
	Timestamp              time.Time
	Service, Metric        string
	Score, Value, Baseline float64
	Reason                 string
}
type Correlation struct {
	Anomaly Anomaly
	Changes []Event
	Score   float64
}
type Hypothesis struct {
	Title, Evidence string
	Probability     float64
	Correlations    []Correlation
}
type Incident struct {
	ID, Service, Severity string
	StartedAt             time.Time
	Hypotheses            []Hypothesis
}
