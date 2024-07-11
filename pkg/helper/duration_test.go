package helper

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testFactor = time.Second

// go test -run TestDurationWithEmptyString
func TestDurationWithEmptyString(t *testing.T) {
	life, err := DurationFromString("", testFactor)

	require.Error(t, err, "Did not return error!")
	require.Equal(t, time.Duration(0), life)
	require.Equal(t, err, ErrIntConvert)
}

// go test -run TestDurationWithInvalidString
func TestDurationWithInvalidString(t *testing.T) {
	life, err := DurationFromString("1m4i", testFactor)

	require.Error(t, err, "Did not return error!")
	require.Equal(t, time.Duration(0), life)
	require.Equal(t, err, ErrIntConvert)
}

// go api -run TestDurationFromString
func TestDurationFromString(t *testing.T) {
	for i := 0; i < 99999; i++ {
		for _, factor := range []time.Duration{time.Nanosecond, time.Microsecond, time.Millisecond, time.Second, time.Minute, time.Hour} {
			life, err := DurationFromString(strconv.Itoa(i+1), factor)

			require.NoError(t, err, "Returned error!")
			require.Equal(t, time.Duration(i+1)*factor, life)
		}
	}
}
