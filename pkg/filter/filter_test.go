package filter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func prepareDatabase() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DryRun: true,
	})
	// Add any setup or seed for the database if needed
	return db
}

type Result struct {
	Name string
}

func TestApplySearchLike(t *testing.T) {
	db := prepareDatabase()

	tests := []struct {
		name     string
		filter   Filter
		columns  []string
		expected string
	}{
		{"NoColumns", Filter{Search: "test"}, []string{}, ""},
		{"NoSearchTerm", Filter{}, []string{"name"}, ""},
		{"SingleColumn", Filter{Search: "test"}, []string{"name"}, "unaccent(LOWER(name)) LIKE unaccent(LOWER('%test%'))"},
		{"MultipleColumns", Filter{Search: "test"}, []string{"name", "description"},
			"unaccent(LOWER(name)) LIKE unaccent(LOWER('%test%')) or unaccent(LOWER(description)) LIKE unaccent(LOWER('%test%'))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ApplySearchLike(db, tt.columns...).Find(&[]Result{}).Statement.SQL.String()
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestApplyOrder(t *testing.T) {
	db := prepareDatabase()

	tests := []struct {
		name      string
		filter    Filter
		tableName *string
		expected  string
	}{
		{"WithoutTable", Filter{Sort: "name", Order: "asc"}, nil, "name asc"},
		{"WithTableName", Filter{Sort: "name", Order: "asc"}, getStringPointer("users"), "users.name asc"},
		{"InvalidOrder", Filter{Sort: "name"}, nil, "name desc"},
		{"InvalidSort", Filter{Order: "asc"}, nil, "updated_at asc"},
	}
	os.Setenv("API_DEFAULT_ORDER", "desc")
	os.Setenv("API_DEFAULT_SORT", "updated_at")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ApplyOrder(db, tt.tableName).Find(&[]Result{}).Statement.SQL.String()
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestApplyPagination(t *testing.T) {
	db := prepareDatabase()

	tests := []struct {
		name     string
		filter   Filter
		expected string
	}{
		{"ValidPageAndLimit", Filter{Page: 2, Limit: 10}, "LIMIT 10 OFFSET 10"},
		{"ZeroPageAndLimit", Filter{}, ""},
		{"NegativePage", Filter{Page: -1, Limit: 10}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ApplyPagination(db).Find(&[]Result{}).Statement.SQL.String()
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestCalcPages(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		count    int64
		expected int64
	}{
		{"ZeroCount", Filter{Limit: 10}, 0, 0},
		{"NormalOperation", Filter{Limit: 10, Page: 2}, 25, 3},
		{"ZeroLimit", Filter{Page: 1}, 25, 1},
		{"ZeroPage", Filter{Limit: 10}, 25, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.CalcPages(tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNew(t *testing.T) {
	sort := "name"
	order := "asc"
	filter := New(sort, order)

	assert.Equal(t, sort, filter.Sort)
	assert.Equal(t, order, filter.Order)
}

func getStringPointer(s string) *string {
	return &s
}
