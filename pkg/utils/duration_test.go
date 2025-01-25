package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationFromString(t *testing.T) {
	tests := []struct {
		str     string
		factor  time.Duration
		want    time.Duration
		wantErr error
	}{
		{"100", time.Second, 100 * time.Second, nil},
		{"0", time.Second, 0, nil},
		{"-50", time.Millisecond, -50 * time.Millisecond, nil},
		{"abc", time.Second, 0, ErrIntConvert},
		{"123abc", time.Second, 0, ErrIntConvert},
		{"", time.Minute, 0, ErrIntConvert},
	}

	for _, tt := range tests {
		got, err := DurationFromString(tt.str, tt.factor)
		assert.Equal(t, tt.want, got)
		assert.Equal(t, tt.wantErr, err)
	}
}
