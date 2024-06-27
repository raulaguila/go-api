package dto

type (
	IDsInputDTO struct {
		IDs []uint `json:"ids"`
	}

	ProfileInputDTO struct {
		Name        *string                `json:"name" example:"ADMIN"`
		Permissions map[string]interface{} `json:"permissions"`
	}

	UserInputDTO struct {
		Name      *string `json:"name" example:"John Cena"`
		Email     *string `json:"email" example:"john.cena@email.com"`
		Status    *bool   `json:"status" example:"true"`
		ProfileID *uint   `json:"profile_id" example:"1"`
	}

	PasswordInputDTO struct {
		Password        *string `json:"password" example:"secret"`
		PasswordConfirm *string `json:"password_confirm" example:"secret"`
	}

	AuthInputDTO struct {
		Login    string `json:"login" example:"admin@admin.com"`
		Password string `json:"password" example:"12345678"`
	}

	DepartmentInputDTO struct {
		Name *string `json:"name" example:"Automation"`
	}
)
