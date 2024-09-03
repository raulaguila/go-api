package mocks

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

type DepartmentServiceMock struct {
	mock.Mock
}

func (s *DepartmentServiceMock) GenerateDepartmentOutputDTO(department *domain.Department) *dto.DepartmentOutputDTO {
	return &dto.DepartmentOutputDTO{
		ID:   &department.ID,
		Name: &department.Name,
	}
}

func (s *DepartmentServiceMock) GetDepartmentByID(ctx context.Context, departmentID uint) (*dto.DepartmentOutputDTO, error) {
	return nil, nil
}

//GetDepartments(context.Context, *filter.Filter) (*dto.ItemsOutputDTO[dto.DepartmentOutputDTO], error)
//CreateDepartment(context.Context, *dto.DepartmentInputDTO) (*dto.DepartmentOutputDTO, error)
//UpdateDepartment(context.Context, uint, *dto.DepartmentInputDTO) (*dto.DepartmentOutputDTO, error)
//DeleteDepartments(context.Context, []uint) error
