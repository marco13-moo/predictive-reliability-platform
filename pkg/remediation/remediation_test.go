package remediation

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"testing"
)

func TestGate(t *testing.T) {
	d := Gate{}.Decide(contracts.Incident{Service: "x", Severity: "high"}, Action{RequiresApproval: true}, false)
	if d.Approved {
		t.Fatal("approved")
	}
}
