package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/stretchr/testify/suite"
)

var (
	ErrItemNotFound = errors.New("item not found")
	ErrInvalidValue = errors.New("invalid value")
)

func TestDepartmentSuit(t *testing.T) {
	suite.Run(t, new(DepartmentTestSuite))
}

type DepartmentTestSuite struct {
	suite.Suite
	ctx    context.Context
	filter *filter.Filter

	service  domain.DepartmentService
	items    []domain.Department
	newItems []domain.Department
	dtos     []dto.DepartmentInputDTO
}

func (s *DepartmentTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.filter = filter.New("name", "desc")
	s.items = []domain.Department{
		{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 1"},
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 2"},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 3"},
		{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 4"},
		{Base: domain.Base{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 5"},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 6"},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "."},
	}
	s.newItems = []domain.Department{
		{Name: "Department 4"},
		{Name: "Department 5"},
		{Name: "."},
	}
	s.dtos = []dto.DepartmentInputDTO{
		{Name: &s.newItems[0].Name},
		{Name: &s.newItems[1].Name},
		{Name: &s.newItems[2].Name},
	}

	var nilFilter *filter.Filter = nil
	var nilDepartment *domain.Department = nil

	repo := &mocks.DepartmentRepositoryMock{}
	repo.On("GetDepartments", s.ctx, s.filter).Return(&s.items, nil)
	repo.On("GetDepartments", s.ctx, nilFilter).Return(&s.items, nil)

	repo.On("CountDepartments", s.ctx, s.filter).Return(int64(len(s.items)), nil)
	repo.On("CountDepartments", s.ctx, nilFilter).Return(int64(0), errors.New("error to count items"))

	repo.On("GetDepartmentByID", s.ctx, s.items[0].ID).Return(&s.items[0], nil)
	repo.On("GetDepartmentByID", s.ctx, s.items[1].ID).Return(&s.items[1], nil)
	repo.On("GetDepartmentByID", s.ctx, s.items[3].ID).Return(&s.items[3], nil)
	repo.On("GetDepartmentByID", s.ctx, s.items[4].ID).Return(&s.items[4], nil)
	repo.On("GetDepartmentByID", s.ctx, s.items[5].ID).Return(&s.items[5], nil)
	repo.On("GetDepartmentByID", s.ctx, uint(10)).Return(nil, ErrItemNotFound)

	repo.On("CreateDepartment", s.ctx, &s.newItems[0]).Return(nil)
	repo.On("CreateDepartment", s.ctx, &s.newItems[1]).Return(nil)
	repo.On("CreateDepartment", s.ctx, nilDepartment).Return(nil, errors.New("error to create item"))

	repo.On("UpdateDepartment", s.ctx, &s.items[3]).Return(nil)
	repo.On("UpdateDepartment", s.ctx, &s.items[4]).Return(nil)

	repo.On("DeleteDepartments", s.ctx, []uint{}).Return(nil)
	repo.On("DeleteDepartments", s.ctx, []uint{1}).Return(nil)
	repo.On("DeleteDepartments", s.ctx, []uint{10}).Return(ErrItemNotFound)

	s.service = NewDepartmentService(repo)
}

func (s *DepartmentTestSuite) TestGetDepartments() {
	items, err := s.service.GetDepartments(s.ctx, s.filter)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.DepartmentOutputDTO]{}, items)
	s.Equal(items.Pagination.TotalItems, uint(len(s.items)))

	_, err = s.service.GetDepartments(s.ctx, nil)

	s.Error(err)
}

func (s *DepartmentTestSuite) TestGetDepartmentByID() {
	item, err := s.service.GetDepartmentByID(s.ctx, s.items[0].ID)

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(s.items[0].Name, *item.Name)
	s.Equal(s.items[0].ID, *item.ID)

	item, err = s.service.GetDepartmentByID(s.ctx, s.items[1].ID)

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(s.items[1].Name, *item.Name)
	s.Equal(s.items[1].ID, *item.ID)

	item, err = s.service.GetDepartmentByID(s.ctx, uint(10))

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

func (s *DepartmentTestSuite) TestCreateDepartment() {
	item, err := s.service.CreateDepartment(s.ctx, &s.dtos[0])

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(s.newItems[0].Name, *item.Name)
	s.Equal(s.newItems[0].ID, *item.ID)

	item, err = s.service.CreateDepartment(s.ctx, &s.dtos[1])

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(s.newItems[1].Name, *item.Name)
	s.Equal(s.newItems[1].ID, *item.ID)

	item, err = s.service.CreateDepartment(s.ctx, nil)

	s.Error(err)
	s.Nil(item)
}

func (s *DepartmentTestSuite) TestUpdateDepartment() {
	item, err := s.service.UpdateDepartment(s.ctx, s.items[3].ID, &s.dtos[0])

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(*s.dtos[0].Name, *item.Name)
	s.Equal(s.items[3].ID, *item.ID)

	item, err = s.service.UpdateDepartment(s.ctx, s.items[4].ID, &s.dtos[1])

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(*s.dtos[1].Name, *item.Name)
	s.Equal(s.items[4].ID, *item.ID)

	item, err = s.service.UpdateDepartment(s.ctx, 10, &s.dtos[1])

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)

	item, err = s.service.UpdateDepartment(s.ctx, s.items[5].ID, &s.dtos[2])

	s.Error(err)
	s.Nil(item)
}

func (s *DepartmentTestSuite) TestDeleteDepartments() {
	s.NoError(s.service.DeleteDepartments(s.ctx, []uint{}))
	s.NoError(s.service.DeleteDepartments(s.ctx, []uint{1}))

	err := s.service.DeleteDepartments(s.ctx, []uint{10})

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
}
