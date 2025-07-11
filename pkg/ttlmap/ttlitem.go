package ttlmap

import (
	"math"
	"time"
)

// NewItem creates an item with the specified value and optional expiration.
func newItem(value any, expiration time.Time, expires bool) *ttlItem {
	return &ttlItem{
		value:      value,
		expiration: expiration,
		expires:    expires,
	}
}

type ttlItem struct {
	value      any
	expiration time.Time
	expires    bool
}

// TTL returns the remaining duration until expiration (negative if expired).
func (s *ttlItem) TTL() time.Duration {
	if s.expires {
		return s.expiration.Sub(time.Now())
	}
	return time.Duration(math.MaxInt64)
}

// Expired checks whether the item is already expired.
func (s *ttlItem) Expired() bool {
	if s.expires {
		return s.expiration.Before(time.Now())
	}
	return false
}
