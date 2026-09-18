package slo

import (
	"fmt"
	"sort"
	"time"

	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
)

type Report struct {
	WindowStart, WindowEnd                                             time.Time
	Target, Total, Bad, Allowance, Remaining, RemainingRatio, BurnRate float64
	Empty, Exhausted                                                   bool
	FormulaVersion                                                     string
}

func Calculate(s contracts.SLO, total, bad int) (Report, error) {
	if s.Target < 0 || s.Target > 1 {
		return Report{}, fmt.Errorf("target must be between 0 and 1")
	}
	if total < 0 || bad < 0 || bad > total {
		return Report{}, fmt.Errorf("invalid eligible/bad counts")
	}
	if total == 0 {
		return Report{Target: s.Target, Empty: true, FormulaVersion: "slo/v1"}, nil
	}
	allowance := float64(total) * (1 - s.Target)
	if allowance <= 0 {
		return Report{}, fmt.Errorf("zero or ambiguous allowance")
	}
	remaining := allowance - float64(bad)
	if remaining < 0 {
		remaining = 0
	}
	return Report{Target: s.Target, Total: float64(total), Bad: float64(bad), Allowance: allowance,
		Remaining: remaining, RemainingRatio: remaining / allowance, BurnRate: float64(bad) / allowance,
		Exhausted: float64(bad) >= allowance, FormulaVersion: "slo/v1"}, nil
}

func ErrorBudget(s contracts.SLO, good, total int, elapsed time.Duration) contracts.Budget {
	r, err := Calculate(s, total, total-good)
	if err != nil || total == 0 {
		return contracts.Budget{SLO: s, Total: s.Window}
	}
	if elapsed <= 0 {
		elapsed = s.Window
	}
	bad := total - good
	consumed := (elapsed / time.Duration(total)) * time.Duration(bad)
	remaining := s.Window - consumed
	if remaining < 0 {
		remaining = 0
	}
	return contracts.Budget{SLO: s, Total: s.Window, Consumed: consumed, Remaining: remaining, BurnRate: r.BurnRate, Exhausted: r.Exhausted}
}

func Availability(events []contracts.Event, s contracts.SLO) (contracts.Budget, error) {
	if s.Target < 0 || s.Target > 1 || s.Window <= 0 {
		return contracts.Budget{}, fmt.Errorf("invalid SLO")
	}
	if len(events) == 0 {
		return ErrorBudget(s, 0, 0, s.Window), nil
	}
	ordered := append([]contracts.Event(nil), events...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Timestamp.Equal(ordered[j].Timestamp) {
			return ordered[i].ID < ordered[j].ID
		}
		return ordered[i].Timestamp.Before(ordered[j].Timestamp)
	})
	good := 0
	for _, e := range ordered {
		if e.Labels["status"] == "ok" || e.Value >= 1 {
			good++
		}
	}
	return ErrorBudget(s, good, len(ordered), ordered[len(ordered)-1].Timestamp.Sub(ordered[0].Timestamp)), nil
}
