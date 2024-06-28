package filters

import (
	"github.com/raulaguila/go-api/pkg/filter"
)

type (
	UserFilter struct {
		filter.Filter
		ProfileID uint `query:"profile_id" form:"profile_id" example:"1"`
	}

	IDFilter struct {
		ID uint `query:"id" form:"id" example:"1"`
	}
)
