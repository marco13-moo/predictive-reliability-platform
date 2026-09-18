package main

import (
	"fmt"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/audit"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/correlate"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/detector"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/fixtures"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/incident"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/remediation"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/replay"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/slo"
	"time"
)

func main() {
	pts := fixtures.Points()
	an := detector.Detector{ZThreshold: 2.5}.Detect("payments", "latency", pts)
	corr := []contracts.Correlation{}
	for _, a := range an {
		corr = append(corr, correlate.Changes(a, fixtures.Events(), 2*time.Minute))
	}
	inc := incident.Rank("payments", corr)
	b, _ := slo.Availability([]contracts.Event{{Labels: map[string]string{"status": "ok"}, Timestamp: time.Now()}}, contracts.SLO{Name: "availability", Service: "payments", Target: .99, Window: 30 * 24 * time.Hour})
	d := remediation.Gate{}.Decide(inc, remediation.Action{Name: "rollback", Risk: "high", RequiresApproval: true}, false)
	tr := audit.Trail{}
	_, _ = tr.Append("demo", "incident.created", inc)
	fmt.Printf("anomalies=%d severity=%s hypotheses=%d budget_remaining=%s approved=%v audit_valid=%v f1=%.2f\n", len(an), inc.Severity, len(inc.Hypotheses), b.Remaining, d.Approved, tr.Verify(), replay.Evaluate(detector.Detector{ZThreshold: 2.5}, pts, []time.Time{pts[len(pts)-1].Timestamp}).F1)
}
