package utils

import (
	"errors"
	"strconv"
	"time"
)

// ErrIntConvert indicates an error when a string cannot be converted to an integer, often due to non-numeric characters.
var ErrIntConvert = errors.New("invalid string number")

// DurationFromString converts a string representation of a number to a time.Duration by applying a given factor.
// Returns an error if the string cannot be converted to an integer.
func DurationFromString(str string, factor time.Duration) (time.Duration, error) {
	converted, err := strconv.Atoi(str)
	if err != nil {
		return time.Duration(0), ErrIntConvert
	}

	return time.Duration(converted) * factor, nil
}
