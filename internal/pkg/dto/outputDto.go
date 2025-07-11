package dto

import (
	"github.com/lib/pq"
)

type (
	ProfileOutputDTO struct {
		ID          *uint           `json:"id" example:"1"`
		Name        *string         `json:"name" example:"ADMIN"`
		Permissions *pq.StringArray `json:"permissions,omitempty"`
	}

	UserOutputDTO struct {
		ID       *uint             `json:"id" example:"1"`
		Name     *string           `json:"name" example:"John Cena"`
		Email    *string           `json:"email" example:"john.cena@email.com"`
		Username *string           `json:"corp_id" example:"john.cena"`
		Status   *bool             `json:"status" example:"true"`
		New      *bool             `json:"new,omitempty" example:"true"`
		Profile  *ProfileOutputDTO `json:"profile,omitempty"`
	}

	outputDTO interface {
		ProfileOutputDTO | UserOutputDTO
	}

	PaginationDTO struct {
		CurrentPage uint `json:"current_page"`
		PageSize    uint `json:"page_size"`
		TotalItems  uint `json:"total_items"`
		TotalPages  uint `json:"total_pages"`
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
