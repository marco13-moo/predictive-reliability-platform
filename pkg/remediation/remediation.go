package remediation

import (
	"fmt"
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
)

type Action struct {
	Name             string
	Risk             string
	RequiresApproval bool
}
type Decision struct {
	Action   Action
	Approved bool
	Reason   string
}
type Policy interface {
	Decide(contracts.Incident, Action, bool) Decision
}
type Gate struct{ AutoApproveLowRisk bool }

func (g Gate) Decide(i contracts.Incident, a Action, approved bool) Decision {
	if a.RequiresApproval && !approved {
		return Decision{Action: a, Reason: "approval required"}
	}
	if i.Severity == "high" && !approved {
		return Decision{Action: a, Reason: "high severity requires approval"}
	}
	if a.Risk == "high" && !g.AutoApproveLowRisk && !approved {
		return Decision{Action: a, Reason: "high risk requires approval"}
	}
	return Decision{Action: a, Approved: true, Reason: fmt.Sprintf("policy accepted for %s", i.Service)}
}
