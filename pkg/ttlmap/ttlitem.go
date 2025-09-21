package ttlmap

import (
	"math"
	"time"
)

// NewItem creates an item with the specified value and optional expiration.
func newItem(value any, expiration *time.Duration) *ttlItem {
	return &ttlItem{
		value:      value,
		createdAt:  time.Now(),
		expiration: expiration,
	}
}

type ttlItem struct {
	value      any
	createdAt  time.Time
	expiration *time.Duration
}

// TTL returns the remaining duration until expiration (negative if expired).
func (s *ttlItem) TTL() time.Duration {
	if s.expiration != nil {
		return time.Until(s.createdAt.Add(*s.expiration))
	}
	return time.Duration(math.MaxInt64)
}

// Expired checks whether the item is already expired.
func (s *ttlItem) Expired() bool {
	if s.expiration != nil {
		return s.createdAt.Add(*s.expiration).Before(time.Now())
	}
	return false
}
