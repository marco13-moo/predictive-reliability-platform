package storage

import (
	"context"
	"time"

	"github.com/marco13-moo/predictive-reliability-platform/pkg/contracts"
)

type Authorization struct {
	TenantID, Subject            string
	CanRead, CanWrite, CanDelete bool
}
type RetentionClass string

const (
	Canonical RetentionClass = "canonical"
	Derived   RetentionClass = "derived"
	Audit     RetentionClass = "audit"
	Replay    RetentionClass = "replay"
)

type Record struct {
	Event     contracts.Event
	Retention RetentionClass
	ExpiresAt time.Time
	Lifecycle string
}
type Store interface {
	Put(context.Context, Authorization, Record) error
	List(context.Context, Authorization, time.Time) ([]Record, error)
	Purge(context.Context, Authorization, time.Time) (int, error)
}
