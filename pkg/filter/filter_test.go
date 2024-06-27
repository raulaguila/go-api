package filter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go api -run TestNewFilter
func TestNewFilter(t *testing.T) {
	_ = os.Setenv("API_DEFAULT_SORT", "updated_at")
	_ = os.Setenv("API_DEFAULT_ORDER", "desc")

	filter := New("updated_at", "desc")

	assert.NotNil(t, filter)
	assert.Equal(t, "", filter.Search)
	assert.Equal(t, 0, filter.Limit)
	assert.Equal(t, 0, filter.Page)
	assert.Equal(t, "updated_at", filter.Sort)
	assert.Equal(t, "desc", filter.Order)
}
