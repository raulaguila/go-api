package dto

import "github.com/lib/pq"

type (
	IDInputDTO[T uint | string] struct {
		ID T `json:"id"`
	}

	IDsInputDTO[T uint | string] struct {
		IDs []T `json:"ids"`
	}

	ProfileInputDTO struct {
		Name        *string         `json:"name" example:"ADMIN"`
		Permissions *pq.StringArray `json:"permissions"`
	}

	UserInputDTO struct {
		Name      *string `json:"name" example:"John Cena"`
		Username  *string `json:"username" example:"john.cena"`
		Email     *string `json:"email" example:"john.cena@email.com"`
		Status    *bool   `json:"status" example:"true"`
		ProfileID *uint   `json:"profile_id" example:"1"`
	}

	PasswordInputDTO struct {
		Password        *string `json:"password" example:"secret"`
		PasswordConfirm *string `json:"password_confirm" example:"secret"`
	}

	AuthInputDTO struct {
		Login      string `json:"login" example:"admin"`
		Password   string `json:"password" example:"12345678"`
		Expiration bool   `json:"expiration" example:"true" default:"true"`
	}
)
