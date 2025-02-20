package dto

import (
	"github.com/raulaguila/go-api/pkg/pgfilter"
)

type (
	ProfileFilter struct {
		pgfilter.Filter
		WithPermissions *bool `query:"with_permissions" form:"with_permissions" example:"false"`
		ListRoot        bool  `query:"list_root" form:"list_root" example:"false"`
	}

	UserFilter struct {
		pgfilter.Filter
		ProfileID uint `query:"profile_id" form:"profile_id" example:"1"`
	}

	IDFilter struct {
		ID uint `query:"id" form:"id" minimum:"1" example:"1" binding:"required"`
	}
)
