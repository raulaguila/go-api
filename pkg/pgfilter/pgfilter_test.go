package pgfilter

import (
	"os"
	"testing"
)

func TestApplySearchLike(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		columns  []string
		expected string
	}{
		{"no_columns", Filter{Search: "test"}, nil, ""},
		{"empty_search", Filter{Search: ""}, []string{"name"}, ""},
		{"single_column", Filter{Search: "test"}, []string{"name"}, "unaccent(LOWER(name)) LIKE unaccent(LOWER('%test%'))"},
		{"multiple_columns", Filter{Search: "test"}, []string{"name", "description"}, "unaccent(LOWER(name)) LIKE unaccent(LOWER('%test%')) or unaccent(LOWER(description)) LIKE unaccent(LOWER('%test%'))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.ApplySearchLike(tt.columns...)
			if got != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestApplyOrder(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		tbName   *string
		expected string
	}{
		{"no_table_name", Filter{Sort: "created_at", Order: "asc"}, nil, "created_at asc"},
		{"table_name_set", Filter{Sort: "created_at", Order: "desc"}, stringPointer("users"), "users.created_at desc"},
		{"no_sort_column", Filter{Sort: "", Order: "asc"}, nil, " asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.ApplyOrder(tt.tbName)
			if got != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestApplyPagination(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		expected bool
		offset   int
		limit    int
	}{
		{"valid_pagination", Filter{Limit: 10, Page: 2}, true, 10, 10},
		{"valid_pagination_2", Filter{Limit: 20, Page: 3}, true, 40, 20},
		{"zero_limit", Filter{Limit: 0, Page: 2}, false, 0, 0},
		{"zero_page", Filter{Limit: 10, Page: 0}, false, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aux, offset, limit := tt.filter.ApplyPagination()
			if aux != tt.expected || offset != tt.offset || limit != tt.limit {
				t.Errorf("expected (%v, %v, %v), got (%v, %v, %v)", tt.expected, tt.offset, tt.limit, aux, offset, limit)
			}
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
		{"zero_count", Filter{Limit: 10}, 0, 0},
		{"count_exact_one_page", Filter{Limit: 10}, 10, 1},
		{"count_multiple_pages", Filter{Limit: 10, Page: 1}, 25, 3},
		{"unlimited_limit", Filter{Limit: 0}, 50, 1},
		{"zero_page", Filter{Limit: 10, Page: 0}, 50, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.CalcPages(tt.count)
			if got != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	defaultOrder := "desc"
	defaultSort := "created_at"

	// Mock environment variables
	os.Setenv("API_DEFAULT_ORDER", defaultOrder)
	os.Setenv("API_DEFAULT_SORT", defaultSort)

	tests := []struct {
		name     string
		filter   Filter
		expected Filter
	}{
		{"valid_order_sort", Filter{Order: "asc", Sort: "updated_at"}, Filter{Order: "asc", Sort: "updated_at"}},
		{"invalid_order", Filter{Order: "wrong", Sort: "created_at"}, Filter{Order: defaultOrder, Sort: "created_at"}},
		{"missing_sort", Filter{Order: "asc", Sort: ""}, Filter{Order: "asc", Sort: defaultSort}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.filter.check()
			if tt.filter.Order != tt.expected.Order || tt.filter.Sort != tt.expected.Sort {
				t.Errorf("expected: %v, got: %v", tt.expected, tt.filter)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		order    string
		expected Filter
	}{
		{"new_filter", "created_at", "asc", Filter{Sort: "created_at", Order: "asc"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.sort, tt.order)
			if got.Sort != tt.expected.Sort || got.Order != tt.expected.Order {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}
