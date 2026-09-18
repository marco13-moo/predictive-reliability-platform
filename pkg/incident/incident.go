package incident

import (
	"fmt"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"sort"
	"time"
)

func Rank(service string, correlations []contracts.Correlation) contracts.Incident {
	sort.Slice(correlations, func(i, j int) bool { return correlations[i].Score > correlations[j].Score })
	hs := make([]contracts.Hypothesis, 0, len(correlations))
	for _, c := range correlations {
		p := c.Score
		if p == 0 {
			p = .1
		}
		hs = append(hs, contracts.Hypothesis{Title: "Recent deployment or configuration change", Evidence: fmt.Sprintf("%d change event(s) near %s", len(c.Changes), c.Anomaly.Timestamp.Format(time.RFC3339)), Probability: p, Correlations: []contracts.Correlation{c}})
	}
	sev := "low"
	if len(hs) > 0 && hs[0].Probability > .5 {
		sev = "high"
	}
	return contracts.Incident{ID: fmt.Sprintf("inc-%d", time.Now().UnixNano()), Service: service, Severity: sev, StartedAt: time.Now().UTC(), Hypotheses: hs}
}
