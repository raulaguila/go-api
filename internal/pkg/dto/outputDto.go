package dto

import (
	"github.com/lib/pq"
)

type (
	ItemOutputDTO struct {
		ID   *uint   `json:"id,omitempty" example:"1"`
		Name *string `json:"name,omitempty" example:"Item"`
	}

	ProfileOutputDTO struct {
		ID          *uint           `json:"id,omitempty" example:"1"`
		Name        *string         `json:"name,omitempty" example:"ADMIN"`
		Permissions *pq.StringArray `json:"permissions,omitempty"`
	}

	UserOutputDTO struct {
		ID       *uint             `json:"id,omitempty" example:"1"`
		Name     *string           `json:"name,omitempty" example:"John Cena"`
		Email    *string           `json:"email,omitempty" example:"john.cena@email.com"`
		Username *string           `json:"corp_id,omitempty" example:"john.cena"`
		Status   *bool             `json:"status,omitempty" example:"true"`
		New      *bool             `json:"new,omitempty" example:"true"`
		Profile  *ProfileOutputDTO `json:"profile,omitempty"`
	}

	outputDTO interface {
		ProfileOutputDTO | UserOutputDTO
	}

	PaginationDTO struct {
		Page       uint `json:"page"`
		Limit      uint `json:"limit"`
		TotalItems uint `json:"total_items"`
		TotalPages uint `json:"total_pages"`
	}

	ItemsOutputDTO[T outputDTO] struct {
		Items      []T           `json:"items"`
		Pagination PaginationDTO `json:"pagination"`
	}

	AuthOutputDTO struct {
		User         *UserOutputDTO `json:"user,omitempty"`
		AccessToken  string         `json:"accesstoken"`
		RefreshToken string         `json:"refreshtoken"`
	}
)
