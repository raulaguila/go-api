package mocks

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/stretchr/testify/mock"
)

type DepartmentRepositoryMock struct {
	mock.Mock
}

func (s *DepartmentRepositoryMock) CountDepartments(ctx context.Context, f *filter.Filter) (int64, error) {
	ret := s.Called(ctx, f)
	return ret.Get(0).(int64), ret.Error(1)
}

func (s *DepartmentRepositoryMock) GetDepartmentByID(ctx context.Context, departmentID uint) (*domain.Department, error) {
	ret := s.Called(ctx, departmentID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.Department), ret.Error(1)
}

func (s *DepartmentRepositoryMock) GetDepartments(ctx context.Context, f *filter.Filter) (*[]domain.Department, error) {
	ret := s.Called(ctx, f)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*[]domain.Department), ret.Error(1)
}

func (s *DepartmentRepositoryMock) CreateDepartment(ctx context.Context, data *dto.DepartmentInputDTO) (*domain.Department, error) {
	ret := s.Called(ctx, data)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.Department), ret.Error(1)
}

func (s *DepartmentRepositoryMock) UpdateDepartment(ctx context.Context, department *domain.Department, data *dto.DepartmentInputDTO) error {
	ret := s.Called(ctx, department, data)
	return ret.Error(0)
}

func (s *DepartmentRepositoryMock) DeleteDepartments(ctx context.Context, departmentIDs []uint) error {
	ret := s.Called(ctx, departmentIDs)
	return ret.Error(0)
}
