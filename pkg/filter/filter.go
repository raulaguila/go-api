package filter

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	"gorm.io/gorm"
)

// orders is a slice containing the strings "asc" and "desc" representing valid order directions for sorting.
var orders = []string{"asc", "desc"}

// New creates a new Filter instance with the specified sort and order parameters.
// The Search, Page, and Limit fields are initialized to their zero values.
func New(sort, order string) *Filter {
	return &Filter{
		Search: "",
		Page:   0,
		Limit:  0,
		Sort:   sort,
		Order:  order,
	}
}

// Filter represents a structure for query parameters used to filter, sort, and paginate API responses.
type Filter struct {
	Search string `query:"search" form:"search" example:"name"`
	Page   int    `query:"page" form:"page" minimum:"1" default:"1"`
	Limit  int    `query:"limit" form:"limit" minimum:"1" default:"10"`
	Sort   string `query:"sort" form:"sort" default:"updated_at" example:"'updated_at', 'created_at', 'name' or some other field from response object"`
	Order  string `query:"order" form:"order" enums:"asc,desc" default:"desc"`
}

// ApplySearchLike applies a search filter on specified columns using a case-insensitive and accent-insensitive matching.
// It returns the modified *gorm.DB pointer that includes the applied filtering condition.
// The method constructs a "WHERE" clause with "LIKE" conditions using the provided column names if `s.Search` is not empty.
// It checks if there are columns specified and if the search term is not an empty string.
// If no conditions are applied, it returns the original database instance.
func (s *Filter) ApplySearchLike(db *gorm.DB, columns ...string) *gorm.DB {
	if len(columns) > 0 && s.Search != "" {
		whereLike := func(column, value string) string {
			return fmt.Sprintf("unaccent(LOWER(%v)) LIKE unaccent(LOWER('%%%v%%'))", column, value)
		}

		where := ""
		for i, column := range columns {
			if i > 0 {
				where += " or "
			}
			where += whereLike(column, s.Search)
		}

		if where != "" {
			return db.Where(where)
		}
	}

	return db
}

// ApplyOrder applies sorting to the database query based on the specified Sort and Order fields in the Filter struct.
// If the optional table name pointer is provided and the Sort field does not contain a dot, the table name is included
// in the ORDER BY clause. The sort order is specified by the Order field, which defaults to the lowercase value.
func (s *Filter) ApplyOrder(db *gorm.DB, tbName *string) *gorm.DB {
	s.check()
	if tbName != nil && !strings.Contains(s.Sort, ".") {
		return db.Order(fmt.Sprintf("%v.%v %v", *tbName, s.Sort, s.Order))
	}
	return db.Order(fmt.Sprintf("%v %v", s.Sort, s.Order))
}

// ApplyPagination applies pagination to the database query by setting the offset and limit based on the current page and limit values of the Filter. If either page or limit is not set, it returns the unmodified database query.
func (s *Filter) ApplyPagination(db *gorm.DB) *gorm.DB {
	if s.Page > 0 && s.Limit > 0 {
		return db.Offset((s.Page - 1) * s.Limit).Limit(s.Limit)
	}

	return db
}

// check validates and normalizes the `Order` and `Sort` fields of the `Filter` struct.
// Converts the `Order` string to lowercase and ensures it matches a valid order, defaulting if invalid.
// Sets the `Sort` field to a default value if it is empty.
func (s *Filter) check() {
	s.Order = strings.ToLower(s.Order)
	if !slices.Contains(orders, s.Order) {
		s.Order = os.Getenv("API_DEFAULT_ORDER")
	}
	if s.Sort == "" {
		s.Sort = os.Getenv("API_DEFAULT_SORT")
	}
}

// CalcPages calculates the total number of pages required to display items.
// It returns 0 if the item count is 0, and 1 if either `Limit` or `Page` is not set.
// Otherwise, it uses ceil division of the item count by the limit.
func (s *Filter) CalcPages(count int64) int64 {
	if count == 0 {
		return 0
	}

	if s.Limit == 0 || s.Page == 0 {
		return 1
	}

	return int64(math.Ceil(float64(count) / float64(s.Limit)))
}
