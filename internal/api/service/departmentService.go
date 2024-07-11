package service

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

func NewDepartmentService(r domain.DepartmentRepository) domain.DepartmentService {
	return &departmentService{
		departmentRepository: r,
	}
}

type departmentService struct {
	departmentRepository domain.DepartmentRepository
}

func (s *departmentService) GenerateDepartmentOutputDTO(department *domain.Department) *dto.DepartmentOutputDTO {
	return &dto.DepartmentOutputDTO{
		ID:   &department.ID,
		Name: &department.Name,
	}
}

// GetDepartmentByID Implementation of 'GetDepartmentByID'.
func (s *departmentService) GetDepartmentByID(ctx context.Context, departmentID uint) (*dto.DepartmentOutputDTO, error) {
	department, err := s.departmentRepository.GetDepartmentByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}

	return s.GenerateDepartmentOutputDTO(department), nil
}

// GetDepartments Implementation of 'GetDepartments'.
func (s *departmentService) GetDepartments(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO[dto.DepartmentOutputDTO], error) {
	departments, err := s.departmentRepository.GetDepartments(ctx, filter)
	if err != nil {
		return nil, err
	}

	count, err := s.departmentRepository.CountDepartments(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputDepartments := make([]dto.DepartmentOutputDTO, 0)
	for _, department := range *departments {
		outputDepartments = append(outputDepartments, *s.GenerateDepartmentOutputDTO(&department))
	}

	return &dto.ItemsOutputDTO[dto.DepartmentOutputDTO]{
		Items: outputDepartments,
		Count: count,
		Pages: filter.CalcPages(count),
	}, nil
}

// CreateDepartment Implementation of 'CreateDepartment'.
func (s *departmentService) CreateDepartment(ctx context.Context, data *dto.DepartmentInputDTO) (*dto.DepartmentOutputDTO, error) {
	department, err := s.departmentRepository.CreateDepartment(ctx, data)
	if err != nil {
		return nil, err
	}

	return s.GenerateDepartmentOutputDTO(department), nil
}

// UpdateDepartment Implementation of 'UpdateDepartment'.
func (s *departmentService) UpdateDepartment(ctx context.Context, departmentID uint, data *dto.DepartmentInputDTO) (*dto.DepartmentOutputDTO, error) {
	department, err := s.departmentRepository.GetDepartmentByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}

	if err := s.departmentRepository.UpdateDepartment(ctx, department, data); err != nil {
		return nil, err
	}

	return s.GenerateDepartmentOutputDTO(department), nil
}

// DeleteDepartments Implementation of 'DeleteDepartments'.
func (s *departmentService) DeleteDepartments(ctx context.Context, ids []uint) error {
	return s.departmentRepository.DeleteDepartments(ctx, ids)
}
