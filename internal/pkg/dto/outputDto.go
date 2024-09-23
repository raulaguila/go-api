package dto

type (
	DepartmentOutputDTO struct {
		ID   *uint   `json:"id,omitempty" example:"1"`
		Name *string `json:"name,omitempty" example:"Automation"`
	}

	ProfileOutputDTO struct {
		ID          *uint          `json:"id" example:"1"`
		Name        *string        `json:"name" example:"ADMIN"`
		Permissions map[string]any `json:"permissions,omitempty"`
	}

	UserOutputDTO struct {
		ID      *uint             `json:"id" example:"1"`
		Name    *string           `json:"name" example:"John Cena"`
		Email   *string           `json:"email" example:"john.cena@email.com"`
		Status  *bool             `json:"status" example:"true"`
		Profile *ProfileOutputDTO `json:"profile,omitempty"`
	}

	outputDTO interface {
		ProfileOutputDTO | UserOutputDTO | DepartmentOutputDTO
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
