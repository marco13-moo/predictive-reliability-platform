package remediation

import (
	"fmt"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"strings"
	"time"
)

type Action struct {
	Name             string
	Risk             string
	RequiresApproval bool
}
type Decision struct {
	Action        Action
	Approved      bool
	Reason        string
	Outcome       string
	ReasonCodes   []string
	PolicyVersion string
	ExpiresAt     time.Time
}
type Policy interface {
	Decide(contracts.Incident, Action, bool) Decision
}
type Gate struct{ AutoApproveLowRisk bool }

func (g Gate) Decide(i contracts.Incident, a Action, approved bool) Decision {
	return g.DecideAt(i, a, approved, time.Now().UTC())
}

func (g Gate) DecideAt(i contracts.Incident, a Action, approved bool, now time.Time) Decision {
	d := Decision{Action: a, PolicyVersion: "remediation/v1", Outcome: "deny"}
	if strings.TrimSpace(i.Service) == "" || strings.TrimSpace(i.ID) == "" {
		d.Reason, d.ReasonCodes = "missing incident scope or evidence", []string{"missing_scope"}
		return d
	}
	if a.RequiresApproval && !approved {
		d.Reason, d.ReasonCodes, d.Outcome = "approval required", []string{"approval_required"}, "require-approval"
		return d
	}
	if i.Severity == "high" && !approved {
		d.Reason, d.ReasonCodes, d.Outcome = "high severity requires approval", []string{"severity_requires_approval"}, "require-approval"
		return d
	}
	if a.Risk == "high" && !g.AutoApproveLowRisk && !approved {
		d.Reason, d.ReasonCodes, d.Outcome = "high risk requires approval", []string{"risk_requires_approval"}, "require-approval"
		return d
	}
	d.Approved, d.Outcome, d.Reason = true, "allow-recommendation", fmt.Sprintf("policy accepted for %s", i.Service)
	d.ExpiresAt = now.UTC().Add(time.Hour)
	return d
}
