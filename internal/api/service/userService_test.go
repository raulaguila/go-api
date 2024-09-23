package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/stretchr/testify/suite"
)

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

type UserTestSuite struct {
	suite.Suite
	ctx context.Context
	f   *filters.UserFilter

	service   domain.UserService
	items     []domain.User
	firstItem domain.User
	newItem   domain.User
}

func (s *UserTestSuite) SetupTest() {
	auth := &domain.Auth{
		Base:      domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Status:    true,
		ProfileID: 1,
		Profile: &domain.Profile{
			Base:        domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:        "ADMIN",
			Permissions: map[string]any{"user": true},
		},
		Token: nil,
	}

	s.ctx = context.Background()
	s.f = &filters.UserFilter{Filter: *filter.New("name", "desc"), ProfileID: 0}
	s.firstItem = domain.User{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 1", Email: "email1@email.com", AuthID: 1, Auth: auth}
	s.items = []domain.User{
		s.firstItem,
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 2", Email: "email2@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 3", Email: "email3@email.com", AuthID: 1, Auth: auth},
	}
	s.newItem = domain.User{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 4", Email: "email4@email.com", AuthID: 1, Auth: auth}

	var nilFilter *filters.UserFilter = nil
	var nilDTO *dto.UserInputDTO = nil

	repo := &mocks.UserRepositoryMock{}
	repo.On("GetUsers", s.ctx, s.f).Return(&s.items, nil)
	repo.On("GetUsers", s.ctx, nilFilter).Return(nil, errors.New("error to get items"))

	repo.On("CountUsers", s.ctx, s.f).Return(int64(len(s.items)), nil)
	repo.On("CountUsers", s.ctx, nilFilter).Return(int64(0), errors.New("error to get items"))

	repo.On("GetUserByID", s.ctx, uint(1)).Return(&s.firstItem, nil)
	repo.On("GetUserByID", s.ctx, uint(4)).Return(&s.newItem, nil)
	repo.On("GetUserByID", s.ctx, uint(7)).Return(nil, ErrItemNotFound)

	repo.On("CreateUser", s.ctx, &dto.UserInputDTO{Name: &s.newItem.Name}).Return(&s.newItem, nil)
	repo.On("CreateUser", s.ctx, nilDTO).Return(nil, errors.New("error to create item"))

	invalidName := "."
	repo.On("UpdateUser", s.ctx, &s.newItem, &dto.UserInputDTO{Name: &s.newItem.Name}).Return(nil)
	repo.On("UpdateUser", s.ctx, &s.newItem, &dto.UserInputDTO{Name: &invalidName}).Return(ErrInvalidValue)

	repo.On("DeleteUsers", s.ctx, []uint{1}).Return(nil)
	repo.On("DeleteUsers", s.ctx, []uint{7}).Return(ErrItemNotFound)

	s.service = NewUserService(repo)
}

func (s *UserTestSuite) TestGetUsers() {
	items, err := s.service.GetUsers(s.ctx, s.f)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.UserOutputDTO]{}, items)
	s.Equal(items.Pagination.TotalItems, uint(len(s.items)))

	items, err = s.service.GetUsers(s.ctx, nil)

	s.Error(err)
	s.Nil(items)
}

func (s *UserTestSuite) TestGetUserByID() {
	item, err := s.service.GetUserByID(s.ctx, 1)

	s.NoError(err)
	s.IsType(&dto.UserOutputDTO{}, item)
	s.Equal(*item.Name, s.firstItem.Name)
	s.Equal(*item.ID, s.firstItem.ID)

	item, err = s.service.GetUserByID(s.ctx, 7)

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

func (s *UserTestSuite) TestCreateUser() {
	data := &dto.UserInputDTO{Name: &s.newItem.Name}
	item, err := s.service.CreateUser(s.ctx, data)

	s.NoError(err)
	s.IsType(&dto.UserOutputDTO{}, item)
	s.Equal(*item.Name, s.newItem.Name)
	s.Equal(*item.ID, s.newItem.ID)

	item, err = s.service.CreateUser(s.ctx, nil)

	s.Error(err)
	s.Equal("error to create item", err.Error())
	s.Nil(item)
}

func (s *UserTestSuite) TestUpdateUser() {
	data := &dto.UserInputDTO{Name: &s.newItem.Name}
	item, err := s.service.UpdateUser(s.ctx, 4, data)

	s.NoError(err)
	s.IsType(&dto.UserOutputDTO{}, item)
	s.Equal(*item.Name, s.newItem.Name)
	s.Equal(*item.ID, s.newItem.ID)

	item, err = s.service.UpdateUser(s.ctx, 7, data)

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)

	invalidName := "."
	data = &dto.UserInputDTO{Name: &invalidName}
	item, err = s.service.UpdateUser(s.ctx, 4, data)

	s.Error(err)
	s.True(errors.Is(err, ErrInvalidValue))
	s.Nil(item)
}

func (s *UserTestSuite) TestDeleteUsers() {
	s.NoError(s.service.DeleteUsers(s.ctx, []uint{1}))

	err := s.service.DeleteUsers(s.ctx, []uint{7})

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
}
