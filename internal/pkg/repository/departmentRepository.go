package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

type departmentRepository struct {
	db *gorm.DB
}

func (s *departmentRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter != nil {
		db = filter.ApplySearchLike(db, "name")
		db = filter.ApplyOrder(db, nil)
	}

	return db
}

func (s *departmentRepository) CountDepartments(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(&domain.Department{}).Count(&count).Error
}

func (s *departmentRepository) GetDepartments(ctx context.Context, filter *filter.Filter) (*[]domain.Department, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	departments := &[]domain.Department{}
	return departments, db.Find(departments).Error
}

func (s *departmentRepository) GetDepartmentByID(ctx context.Context, departmentID uint) (*domain.Department, error) {
	department := &domain.Department{}
	return department, s.db.WithContext(ctx).First(department, departmentID).Error
}

func (s *departmentRepository) CreateDepartment(ctx context.Context, data *dto.DepartmentInputDTO) (*domain.Department, error) {
	department := &domain.Department{}
	if err := department.Bind(data); err != nil {
		return nil, err
	}

	return department, s.db.WithContext(ctx).Create(department).Error
}

func (s *departmentRepository) UpdateDepartment(ctx context.Context, department *domain.Department, data *dto.DepartmentInputDTO) error {
	if err := department.Bind(data); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(department).Updates(department.ToMap()).Error
}

func (s *departmentRepository) DeleteDepartments(ctx context.Context, toDelete []uint) error {
	if len(toDelete) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&domain.Department{}, toDelete).Error
	})
}
