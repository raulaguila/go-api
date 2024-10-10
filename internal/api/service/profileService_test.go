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

func TestProfileSuit(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

type ProfileTestSuite struct {
	suite.Suite
	ctx    context.Context
	filter *filter.Filter

	service  domain.ProfileService
	items    []domain.Profile
	newItems []domain.Profile
	dtos     []dto.ProfileInputDTO
}

func (s *ProfileTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.filter = filter.New("name", "desc")
	s.items = []domain.Profile{
		{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 1", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 2", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 3", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 4", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 5", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 6", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: ".", Permissions: map[string]any{"profile": true}},
	}
	s.newItems = []domain.Profile{
		{Name: "Profile 4", Permissions: map[string]any{"profile": true}},
		{Name: "Profile 5", Permissions: map[string]any{"profile": true}},
		{Name: ".", Permissions: map[string]any{"profile": true}},
	}
	s.dtos = []dto.ProfileInputDTO{
		{Name: &s.newItems[0].Name, Permissions: s.newItems[0].Permissions},
		{Name: &s.newItems[1].Name, Permissions: s.newItems[1].Permissions},
		{Name: &s.newItems[2].Name, Permissions: s.newItems[2].Permissions},
	}

	var nilFilter *filter.Filter = nil
	var nilProfile *domain.Profile = nil

	repo := &mocks.ProfileRepositoryMock{}
	repo.On("GetProfiles", s.ctx, s.filter).Return(&s.items, nil)
	repo.On("GetProfiles", s.ctx, nilFilter).Return(&s.items, nil)

	repo.On("CountProfiles", s.ctx, s.filter).Return(int64(len(s.items)), nil)
	repo.On("CountProfiles", s.ctx, nilFilter).Return(int64(0), errors.New("error to count items"))

	repo.On("GetProfileByID", s.ctx, s.items[0].ID).Return(&s.items[0], nil)
	repo.On("GetProfileByID", s.ctx, s.items[1].ID).Return(&s.items[1], nil)
	repo.On("GetProfileByID", s.ctx, s.items[3].ID).Return(&s.items[3], nil)
	repo.On("GetProfileByID", s.ctx, s.items[4].ID).Return(&s.items[4], nil)
	repo.On("GetProfileByID", s.ctx, s.items[5].ID).Return(&s.items[5], nil)
	repo.On("GetProfileByID", s.ctx, uint(10)).Return(nil, ErrItemNotFound)

	repo.On("CreateProfile", s.ctx, &s.newItems[0]).Return(nil)
	repo.On("CreateProfile", s.ctx, &s.newItems[1]).Return(nil)
	repo.On("CreateProfile", s.ctx, nilProfile).Return(nil, errors.New("error to create item"))

	repo.On("UpdateProfile", s.ctx, &s.items[3]).Return(nil)
	repo.On("UpdateProfile", s.ctx, &s.items[4]).Return(nil)

	repo.On("DeleteProfiles", s.ctx, []uint{}).Return(nil)
	repo.On("DeleteProfiles", s.ctx, []uint{1}).Return(nil)
	repo.On("DeleteProfiles", s.ctx, []uint{10}).Return(ErrItemNotFound)

	s.service = NewProfileService(repo)
}

func (s *ProfileTestSuite) TestGetProfiles() {
	items, err := s.service.GetProfiles(s.ctx, s.filter)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.ProfileOutputDTO]{}, items)
	s.Equal(items.Pagination.TotalItems, uint(len(s.items)))

	_, err = s.service.GetProfiles(s.ctx, nil)

	s.Error(err)
}

func (s *ProfileTestSuite) TestGetProfileByID() {
	item, err := s.service.GetProfileByID(s.ctx, s.items[0].ID)

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(s.items[0].Name, *item.Name)
	s.Equal(s.items[0].ID, *item.ID)

	item, err = s.service.GetProfileByID(s.ctx, s.items[1].ID)

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(s.items[1].Name, *item.Name)
	s.Equal(s.items[1].ID, *item.ID)

	item, err = s.service.GetProfileByID(s.ctx, uint(10))

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

func (s *ProfileTestSuite) TestCreateProfile() {
	item, err := s.service.CreateProfile(s.ctx, &s.dtos[0])

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(s.newItems[0].Name, *item.Name)
	s.Equal(s.newItems[0].ID, *item.ID)

	item, err = s.service.CreateProfile(s.ctx, &s.dtos[1])

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(s.newItems[1].Name, *item.Name)
	s.Equal(s.newItems[1].ID, *item.ID)

	item, err = s.service.CreateProfile(s.ctx, nil)

	s.Error(err)
	s.Nil(item)
}

func (s *ProfileTestSuite) TestUpdateProfile() {
	item, err := s.service.UpdateProfile(s.ctx, s.items[3].ID, &s.dtos[0])

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(*s.dtos[0].Name, *item.Name)
	s.Equal(s.items[3].ID, *item.ID)

	item, err = s.service.UpdateProfile(s.ctx, s.items[4].ID, &s.dtos[1])

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(*s.dtos[1].Name, *item.Name)
	s.Equal(s.items[4].ID, *item.ID)

	item, err = s.service.UpdateProfile(s.ctx, 10, &s.dtos[1])

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)

	item, err = s.service.UpdateProfile(s.ctx, s.items[5].ID, &s.dtos[2])

	s.Error(err)
	s.Nil(item)
}

func (s *ProfileTestSuite) TestDeleteProfiles() {
	s.NoError(s.service.DeleteProfiles(s.ctx, []uint{}))
	s.NoError(s.service.DeleteProfiles(s.ctx, []uint{1}))

	err := s.service.DeleteProfiles(s.ctx, []uint{10})

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
}
