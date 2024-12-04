package helper

import (
	"testing"
	"time"
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
		if got != tt.want || err != tt.wantErr {
			t.Errorf("DurationFromString(%q, %v) = %v, %v; want %v, %v", tt.str, tt.factor, got, err, tt.want, tt.wantErr)
		}
	}
}
