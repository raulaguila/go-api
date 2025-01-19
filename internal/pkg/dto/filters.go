package dto

import (
	"github.com/raulaguila/go-api/pkg/pgfilter"
)

type (
	UserFilter struct {
		pgfilter.Filter
		ProfileID uint `query:"profile_id" form:"profile_id" example:"1"`
	}

	IDFilter struct {
		ID uint `query:"id" form:"id" minimum:"1" example:"1" binding:"required"`
	}
)
