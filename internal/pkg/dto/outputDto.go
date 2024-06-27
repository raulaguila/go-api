package dto

type (
	DepartmentOutputDTO struct {
		Id   *uint   `json:"id,omitempty" example:"1"`
		Name *string `json:"name,omitempty" example:"Automation"`
	}

	ProfileOutputDTO struct {
		Id          *uint                  `json:"id" example:"1"`
		Name        *string                `json:"name" example:"ADMIN"`
		Permissions map[string]interface{} `json:"permissions,omitempty"`
	}

	UserOutputDTO struct {
		Id      *uint             `json:"id" example:"1"`
		Name    *string           `json:"name" example:"John Cena"`
		Email   *string           `json:"email" example:"john.cena@email.com"`
		Status  *bool             `json:"status" example:"true"`
		Profile *ProfileOutputDTO `json:"profile,omitempty"`
	}

	outputDTO interface {
		ProfileOutputDTO | UserOutputDTO | DepartmentOutputDTO
	}

	ItemsOutputDTO[T outputDTO] struct {
		Items []T   `json:"items"`
		Count int64 `json:"count"`
		Pages int64 `json:"pages"`
	}

	AuthOutputDTO struct {
		User         *UserOutputDTO `json:"user,omitempty"`
		AccessToken  string         `json:"accesstoken"`
		RefreshToken string         `json:"refreshtoken"`
	}
)
