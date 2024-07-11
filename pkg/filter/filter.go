package filter

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var orders = []string{"asc", "desc"}

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
	Sort   string `query:"sort" form:"sort" default:"updated_at" example:"'updated_at', 'created_at', 'name' or some other field from response object"`
	Order  string `query:"order" form:"order" enums:"asc,desc" default:"desc"`
}

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

func (s *Filter) ApplyOrder(db *gorm.DB, tbName *string) *gorm.DB {
	s.check()
	if tbName != nil && !strings.Contains(s.Sort, ".") {
		return db.Order(fmt.Sprintf("%v.%v %v", *tbName, s.Sort, s.Order))
	}
	return db.Order(fmt.Sprintf("%v %v", s.Sort, s.Order))
}

func (s *Filter) ApplyPagination(db *gorm.DB) *gorm.DB {
	if s.Page > 0 && s.Limit > 0 {
		return db.Offset((s.Page - 1) * s.Limit).Limit(s.Limit)
	}

	return db
}

func (s *Filter) check() {
	s.Order = strings.ToLower(s.Order)
	if !slices.Contains(orders, s.Order) {
		s.Order = os.Getenv("API_DEFAULT_ORDER")
	}
	if s.Sort == "" {
		s.Sort = os.Getenv("API_DEFAULT_SORT")
	}
}

func (s *Filter) CalcPages(count int64) int64 {
	if s.Limit == 0 {
		return 1
	}

	return int64(math.Ceil(float64(count) / float64(s.Limit)))
}
