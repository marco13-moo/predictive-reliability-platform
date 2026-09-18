package contracts

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type EventKind string

const (
	KindMetric EventKind = "metric"
	KindLog    EventKind = "log"
	KindChange EventKind = "change"
)

const CurrentSchemaVersion = 1

type Provenance struct {
	Connector, Region, ReceivedBy string
	ParserVersion                 string
}
type Event struct {
	ID            string
	SchemaVersion int
	TenantID      string
	Kind          EventKind
	Source        string
	Service       string
	Timestamp     time.Time
	IngestionTime time.Time
	Name          string
	Unit          string
	Semantic      string
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
	if e.SchemaVersion > CurrentSchemaVersion {
		return fmt.Errorf("unsupported schema version %d", e.SchemaVersion)
	}
	if strings.TrimSpace(e.TenantID) == "" {
		return fmt.Errorf("tenant scope is required")
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	if e.Timestamp.Location() != time.UTC {
		return fmt.Errorf("timestamp must be UTC")
	}
	if !e.IngestionTime.IsZero() && e.IngestionTime.Location() != time.UTC {
		return fmt.Errorf("ingestion time must be UTC")
	}
	if e.Kind != KindMetric && e.Kind != KindLog && e.Kind != KindChange {
		return fmt.Errorf("unsupported event kind %q", e.Kind)
	}
	if math.IsNaN(e.Value) || math.IsInf(e.Value, 0) {
		return fmt.Errorf("value must be finite")
	}
	if e.Kind == KindMetric && strings.TrimSpace(e.Unit) == "" {
		return fmt.Errorf("metric unit is required")
	}
	if strings.TrimSpace(e.Provenance.Connector) == "" || strings.TrimSpace(e.Provenance.ParserVersion) == "" {
		return fmt.Errorf("provenance connector and parser version are required")
	}
	if len(e.Labels) > 64 {
		return fmt.Errorf("too many labels")
	}
	for k, v := range e.Labels {
		if strings.TrimSpace(k) == "" || len(k) > 128 || len(v) > 512 {
			return fmt.Errorf("label exceeds bounds")
		}
	}
	return nil
}

type MetricPoint struct {
	ID        string
	Timestamp time.Time
	Value     float64
	TenantID  string
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
	Timestamp                  time.Time
	TenantID, Service, Metric  string
	Score, Value, Baseline     float64
	Reason                     string
	DetectorVersion            string
	Threshold                  float64
	EvidenceStart, EvidenceEnd time.Time
	Provenance                 Provenance
}
type Correlation struct {
	Anomaly Anomaly
	Changes []Event
	Score   float64
	Factors []string
	Window  time.Duration
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
