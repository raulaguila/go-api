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
	ctx    context.Context
	filter *filters.UserFilter

	service  domain.UserService
	items    []domain.User
	newItems []domain.User
	dtos     []dto.UserInputDTO
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
	s.filter = &filters.UserFilter{Filter: *filter.New("name", "desc"), ProfileID: 0}
	s.items = []domain.User{
		{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 1", Email: "email1@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 2", Email: "email2@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 3", Email: "email3@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 4", Email: "email4@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 5", Email: "email5@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "User 6", Email: "email6@email.com", AuthID: 1, Auth: auth},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: ".", Email: ".", AuthID: 1, Auth: auth},
	}
	var um uint = 1
	s.newItems = []domain.User{
		{Name: "User 4", Email: "email4@email.com", AuthID: 0, Auth: &domain.Auth{ProfileID: um, Profile: auth.Profile}},
		{Name: "User 5", Email: "email5@email.com", AuthID: 0, Auth: &domain.Auth{ProfileID: um, Profile: auth.Profile}},
		{Name: ".", Email: ".", AuthID: 0, Auth: &domain.Auth{ProfileID: um, Profile: auth.Profile}},
	}
	s.dtos = []dto.UserInputDTO{
		{Name: &s.newItems[0].Name, Email: &s.newItems[0].Email, ProfileID: &um},
		{Name: &s.newItems[1].Name, Email: &s.newItems[1].Email, ProfileID: &um},
		{Name: &s.newItems[2].Name, Email: &s.newItems[2].Email, ProfileID: &um},
	}

	var nilFilter *filters.UserFilter = nil
	// var nilUser *domain.User = nil

	repo := &mocks.UserRepositoryMock{}
	repo.On("GetUsers", s.ctx, s.filter).Return(&s.items, nil)
	repo.On("GetUsers", s.ctx, nilFilter).Return(&s.items, nil)

	repo.On("CountUsers", s.ctx, s.filter).Return(int64(len(s.items)), nil)
	repo.On("CountUsers", s.ctx, nilFilter).Return(int64(0), errors.New("error to count items"))

	repo.On("GetUserByID", s.ctx, s.items[0].ID).Return(&s.items[0], nil)
	repo.On("GetUserByID", s.ctx, s.items[1].ID).Return(&s.items[1], nil)
	repo.On("GetUserByID", s.ctx, s.items[3].ID).Return(&s.items[3], nil)
	repo.On("GetUserByID", s.ctx, s.items[4].ID).Return(&s.items[4], nil)
	repo.On("GetUserByID", s.ctx, s.items[5].ID).Return(&s.items[5], nil)
	repo.On("GetUserByID", s.ctx, uint(10)).Return(nil, ErrItemNotFound)

	// repo.On("CreateUser", s.ctx, &s.newItems[0]).Return(nil)
	// repo.On("CreateUser", s.ctx, &s.newItems[1]).Return(nil)
	// repo.On("CreateUser", s.ctx, nilUser).Return(nil, errors.New("error to create item"))

	//repo.On("UpdateUser", s.ctx, &s.newItem).Return(nil)
	//s.newItem.Name = "."
	//repo.On("UpdateUser", s.ctx, &s.newItem).Return(ErrInvalidValue)
	//
	//repo.On("DeleteUsers", s.ctx, []uint{1}).Return(nil)
	//repo.On("DeleteUsers", s.ctx, []uint{7}).Return(ErrItemNotFound)

	s.service = NewUserService(repo)
}

func (s *UserTestSuite) TestGetUsers() {
	items, err := s.service.GetUsers(s.ctx, s.filter)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.UserOutputDTO]{}, items)
	s.Equal(items.Pagination.TotalItems, uint(len(s.items)))

	_, err = s.service.GetUsers(s.ctx, nil)

	s.Error(err)
}

func (s *UserTestSuite) TestGetUserByID() {
	item, err := s.service.GetUserByID(s.ctx, s.items[0].ID)

	s.NoError(err)
	s.IsType(&dto.UserOutputDTO{}, item)
	s.Equal(s.items[0].Name, *item.Name)
	s.Equal(s.items[0].ID, *item.ID)

	item, err = s.service.GetUserByID(s.ctx, s.items[1].ID)

	s.NoError(err)
	s.IsType(&dto.UserOutputDTO{}, item)
	s.Equal(s.items[1].Name, *item.Name)
	s.Equal(s.items[1].ID, *item.ID)

	item, err = s.service.GetUserByID(s.ctx, uint(10))

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

//func (s *UserTestSuite) TestCreateUser() {
//	item, err := s.service.CreateUser(s.ctx, &s.dtos[0])
//
//	s.NoError(err)
//	s.IsType(&dto.UserOutputDTO{}, item)
//	s.Equal(s.newItems[0].Name, *item.Name)
//	s.Equal(s.newItems[0].ID, *item.ID)
//
//	item, err = s.service.CreateUser(s.ctx, &s.dtos[1])
//
//	s.NoError(err)
//	s.IsType(&dto.UserOutputDTO{}, item)
//	s.Equal(s.newItems[1].Name, *item.Name)
//	s.Equal(s.newItems[1].ID, *item.ID)
//
//	item, err = s.service.CreateUser(s.ctx, nil)
//
//	s.Error(err)
//	s.Nil(item)
//}

//func (s *UserTestSuite) TestUpdateUser() {
//	data := &dto.UserInputDTO{Name: &s.newItem.Name}
//	item, err := s.service.UpdateUser(s.ctx, 4, data)
//
//	s.NoError(err)
//	s.IsType(&dto.UserOutputDTO{}, item)
//	s.Equal(*item.Name, s.newItem.Name)
//	s.Equal(*item.ID, s.newItem.ID)
//
//	item, err = s.service.UpdateUser(s.ctx, 7, data)
//
//	s.Error(err)
//	s.True(errors.Is(err, ErrItemNotFound))
//	s.Nil(item)
//
//	invalidName := "."
//	data = &dto.UserInputDTO{Name: &invalidName}
//	item, err = s.service.UpdateUser(s.ctx, 4, data)
//
//	s.Error(err)
//	s.True(errors.Is(err, ErrInvalidValue))
//	s.Nil(item)
//}
//
//func (s *UserTestSuite) TestDeleteUsers() {
//	s.NoError(s.service.DeleteUsers(s.ctx, []uint{1}))
//
//	err := s.service.DeleteUsers(s.ctx, []uint{7})
//
//	s.Error(err)
//	s.True(errors.Is(err, ErrItemNotFound))
//}
