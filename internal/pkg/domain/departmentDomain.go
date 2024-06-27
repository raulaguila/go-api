package domain

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/validator"
)

const DepartmentTableName string = "department"

type (
	Department struct {
		Base
		Name string `gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	DepartmentRepository interface {
		CountDepartments(context.Context, *filter.Filter) (int64, error)
		GetDepartmentByID(context.Context, uint) (*Department, error)
		GetDepartments(context.Context, *filter.Filter) (*[]Department, error)
		CreateDepartment(context.Context, *dto.DepartmentInputDTO) (*Department, error)
		UpdateDepartment(context.Context, *Department, *dto.DepartmentInputDTO) error
		DeleteDepartments(context.Context, []uint) error
	}

	DepartmentService interface {
		GenerateDepartmentOutputDTO(*Department) *dto.DepartmentOutputDTO
		GetDepartmentByID(context.Context, uint) (*dto.DepartmentOutputDTO, error)
		GetDepartments(context.Context, *filter.Filter) (*dto.ItemsOutputDTO[dto.DepartmentOutputDTO], error)
		CreateDepartment(context.Context, *dto.DepartmentInputDTO) (*dto.DepartmentOutputDTO, error)
		UpdateDepartment(context.Context, uint, *dto.DepartmentInputDTO) (*dto.DepartmentOutputDTO, error)
		DeleteDepartments(context.Context, []uint) error
	}
)

func (s *Department) TableName() string { return DepartmentTableName }

func (s *Department) Bind(p *dto.DepartmentInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	return validator.StructValidator.Validate(s)
}

func (s *Department) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": s.Name,
	}
}
