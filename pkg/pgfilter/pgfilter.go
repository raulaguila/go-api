package pgfilter

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)

func New(sort, order string) *Filter {
	return &Filter{
		Search: "",
		Page:   0,
		Limit:  0,
		Sort:   sort,
		Order:  order,
	}
}

type Filter struct {
	Search string `query:"search" form:"search" example:"name"`
	Page   int    `query:"page" form:"page" minimum:"1" default:"1"`
	Limit  int    `query:"limit" form:"limit" minimum:"1" default:"10"`
	Sort   string `query:"sort" form:"sort" default:"updated_at" example:"updated_at"`
	Order  string `query:"order" form:"order" enums:"asc,desc" default:"desc"`
}

func (s *Filter) ApplySearchLike(columns ...string) (where string) {
	if len(columns) > 0 && s.Search != "" {
		whereFunc := func(column, value string) string {
			return fmt.Sprintf("unaccent(LOWER(%s)) LIKE unaccent(LOWER('%%%s%%'))", column, value)
		}

		for i, column := range columns {
			if i > 0 {
				where += " or "
			}
			where += whereFunc(column, s.Search)
		}
	}

	return
}

func (s *Filter) ApplyOrder(tbName *string) string {
	s.check()
	if tbName != nil && !strings.Contains(s.Sort, ".") {
		return fmt.Sprintf("%v.%v %v", *tbName, s.Sort, s.Order)
	}
	return fmt.Sprintf("%v %v", s.Sort, s.Order)
}

func (s *Filter) ApplyPagination() (aux bool, offset, limit int) {
	if aux = s.Page > 0 && s.Limit > 0; aux {
		offset = (s.Page - 1) * s.Limit
		limit = s.Limit
	}

	return
}

func (s *Filter) check() {
	if !slices.Contains([]string{"asc", "desc"}, strings.ToLower(s.Order)) {
		s.Order = os.Getenv("API_DEFAULT_ORDER")
	}
	if s.Sort == "" {
		s.Sort = os.Getenv("API_DEFAULT_SORT")
	}
}

func (s *Filter) CalcPages(count int64) int64 {
	if count == 0 {
		return 0
	}

	if s.Limit == 0 || s.Page == 0 {
		return 1
	}

	return int64(math.Ceil(float64(count) / float64(s.Limit)))
}
