package dto

import (
	"github.com/raulaguila/go-api/pkg/pgfilter"
)

type (
	IDFilter[T uint | string] struct {
		ID T `query:"id" form:"id" minimum:"1" example:"1" binding:"required"`
	}

	IDsFilter[T uint | string] struct {
		IDs []T `query:"ids" form:"ids" minimum:"1" example:"1" binding:"required"`
	}

	ProfileFilter struct {
		pgfilter.Filter
		WithPermissions *bool `query:"with_permissions" form:"with_permissions" example:"false"`
		ListRoot        bool  `query:"list_root" form:"list_root" example:"false"`
	}

	UserFilter struct {
		pgfilter.Filter
		ProfileID uint  `query:"profile_id" form:"level_id" example:"1"`
		Status    *bool `query:"status" form:"status" example:"false"`
	}

	EvidenceFilter struct {
		pgfilter.Filter
		EmployeeID uint `query:"employee_id" form:"employee_id" example:"1"`
	}
)
