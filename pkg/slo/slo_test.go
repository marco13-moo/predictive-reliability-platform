package slo

import (
	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
	"testing"
	"time"
)

func TestErrorBudgetMath(t *testing.T) {
	b := ErrorBudget(contracts.SLO{Target: .99, Window: 100 * time.Hour}, 99, 100, 100*time.Hour)
	if b.Consumed != time.Hour || b.Remaining != 99*time.Hour || b.BurnRate < .999 || b.BurnRate > 1.001 {
		t.Fatalf("%+v", b)
	}
}
