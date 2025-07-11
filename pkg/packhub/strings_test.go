package packhub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		Name   string
		Input  string
		Output string
	}{
		{
			Name:   "OK 1",
			Input:  "normal string",
			Output: "Normal string",
		},
		{
			Name:   "OK 2",
			Input:  "NORMAL STRING",
			Output: "Normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Output, Capitalize(tt.Input))
		})
	}
}
