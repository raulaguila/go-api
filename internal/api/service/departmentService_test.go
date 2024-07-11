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
	ctx context.Context
	f   *filter.Filter

	service   domain.DepartmentService
	items     []domain.Department
	firstItem domain.Department
	newItem   domain.Department
}

func (s *DepartmentTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.f = filter.New("name", "desc")
	s.firstItem = domain.Department{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 1"}
	s.items = []domain.Department{
		s.firstItem,
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 2"},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 3"},
	}
	s.newItem = domain.Department{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Department 4"}

	var nilFilter *filter.Filter = nil
	var nilDTO *dto.DepartmentInputDTO = nil

	repo := &mocks.DepartmentRepositoryMock{}
	repo.On("GetDepartments", s.ctx, s.f).Return(&s.items, nil)
	repo.On("GetDepartments", s.ctx, nilFilter).Return(nil, errors.New("error to get items"))

	repo.On("CountDepartments", s.ctx, s.f).Return(int64(len(s.items)), nil)
	repo.On("CountDepartments", s.ctx, nilFilter).Return(int64(0), errors.New("error to get items"))

	repo.On("GetDepartmentByID", s.ctx, uint(1)).Return(&s.firstItem, nil)
	repo.On("GetDepartmentByID", s.ctx, uint(4)).Return(&s.newItem, nil)
	repo.On("GetDepartmentByID", s.ctx, uint(7)).Return(nil, ErrItemNotFound)

	repo.On("CreateDepartment", s.ctx, &dto.DepartmentInputDTO{Name: &s.newItem.Name}).Return(&s.newItem, nil)
	repo.On("CreateDepartment", s.ctx, nilDTO).Return(nil, errors.New("error to create item"))

	invalidName := "."
	repo.On("UpdateDepartment", s.ctx, &s.newItem, &dto.DepartmentInputDTO{Name: &s.newItem.Name}).Return(nil)
	repo.On("UpdateDepartment", s.ctx, &s.newItem, &dto.DepartmentInputDTO{Name: &invalidName}).Return(ErrInvalidValue)

	repo.On("DeleteDepartments", s.ctx, []uint{1}).Return(nil)
	repo.On("DeleteDepartments", s.ctx, []uint{7}).Return(ErrItemNotFound)

	s.service = NewDepartmentService(repo)
}

func (s *DepartmentTestSuite) TestGetDepartments() {
	items, err := s.service.GetDepartments(s.ctx, s.f)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.DepartmentOutputDTO]{}, items)
	s.Equal(items.Count, int64(len(s.items)))

	items, err = s.service.GetDepartments(s.ctx, nil)

	s.Error(err)
	s.Nil(items)
}

func (s *DepartmentTestSuite) TestGetDepartmentByID() {
	item, err := s.service.GetDepartmentByID(s.ctx, 1)

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(*item.Name, s.firstItem.Name)
	s.Equal(*item.ID, s.firstItem.ID)

	item, err = s.service.GetDepartmentByID(s.ctx, 7)

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

func (s *DepartmentTestSuite) TestCreateDepartment() {
	data := &dto.DepartmentInputDTO{Name: &s.newItem.Name}
	item, err := s.service.CreateDepartment(s.ctx, data)

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(*item.Name, s.newItem.Name)
	s.Equal(*item.ID, s.newItem.ID)

	item, err = s.service.CreateDepartment(s.ctx, nil)

	s.Error(err)
	s.Equal("error to create item", err.Error())
	s.Nil(item)
}

func (s *DepartmentTestSuite) TestUpdateDepartment() {
	data := &dto.DepartmentInputDTO{Name: &s.newItem.Name}
	item, err := s.service.UpdateDepartment(s.ctx, 4, data)

	s.NoError(err)
	s.IsType(&dto.DepartmentOutputDTO{}, item)
	s.Equal(*item.Name, s.newItem.Name)
	s.Equal(*item.ID, s.newItem.ID)

	item, err = s.service.UpdateDepartment(s.ctx, 7, data)

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)

	invalidName := "."
	data = &dto.DepartmentInputDTO{Name: &invalidName}
	item, err = s.service.UpdateDepartment(s.ctx, 4, data)

	s.Error(err)
	s.True(errors.Is(err, ErrInvalidValue))
	s.Nil(item)
}

func (s *DepartmentTestSuite) TestDeleteDepartments() {
	s.NoError(s.service.DeleteDepartments(s.ctx, []uint{1}))

	err := s.service.DeleteDepartments(s.ctx, []uint{7})

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
}
