package slo

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"sort"
	"time"
)

func ErrorBudget(s contracts.SLO, good, total int, elapsed time.Duration) contracts.Budget {
	allowed := float64(total) * (1 - s.Target)
	bad := float64(total - good)
	consumed := time.Duration(0)
	if total > 0 {
		consumed = time.Duration(float64(elapsed) * bad / float64(total))
	}
	remaining := s.Window - consumed
	if remaining < 0 {
		remaining = 0
	}
	burn := 0.0
	if allowed > 0 && total > 0 {
		burn = bad / (float64(total) * (1 - s.Target))
	}
	return contracts.Budget{SLO: s, Total: s.Window, Consumed: consumed, Remaining: remaining, BurnRate: burn, Exhausted: bad > allowed}
}
func Availability(events []contracts.Event, s contracts.SLO) (contracts.Budget, error) {
	if len(events) == 0 {
		return ErrorBudget(s, 0, 0, s.Window), nil
	}
	sort.Slice(events, func(i, j int) bool { return events[i].Timestamp.Before(events[j].Timestamp) })
	good := 0
	for _, e := range events {
		if e.Labels["status"] == "ok" || e.Value >= 1 {
			good++
		}
	}
	elapsed := events[len(events)-1].Timestamp.Sub(events[0].Timestamp)
	if elapsed <= 0 {
		elapsed = s.Window
	}
	return ErrorBudget(s, good, len(events), elapsed), nil
}
