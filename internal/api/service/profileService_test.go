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
	ctx context.Context
	f   *filter.Filter

	service   domain.ProfileService
	items     []domain.Profile
	firstItem domain.Profile
	newItem   domain.Profile
}

func (s *ProfileTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.f = filter.New("name", "desc")
	s.firstItem = domain.Profile{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 1", Permissions: map[string]any{"profile": true}}
	s.items = []domain.Profile{
		s.firstItem,
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 2", Permissions: map[string]any{"profile": true}},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 3", Permissions: map[string]any{"profile": true}},
	}
	s.newItem = domain.Profile{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Profile 4", Permissions: map[string]any{"profile": true}}

	var nilFilter *filter.Filter = nil
	var nilDTO *dto.ProfileInputDTO = nil

	repo := &mocks.ProfileRepositoryMock{}
	repo.On("GetProfiles", s.ctx, s.f).Return(&s.items, nil)
	repo.On("GetProfiles", s.ctx, nilFilter).Return(nil, errors.New("error to get items"))

	repo.On("CountProfiles", s.ctx, s.f).Return(int64(len(s.items)), nil)
	repo.On("CountProfiles", s.ctx, nilFilter).Return(int64(0), errors.New("error to get items"))

	repo.On("GetProfileByID", s.ctx, uint(1)).Return(&s.firstItem, nil)
	repo.On("GetProfileByID", s.ctx, uint(4)).Return(&s.newItem, nil)
	repo.On("GetProfileByID", s.ctx, uint(7)).Return(nil, ErrItemNotFound)

	repo.On("CreateProfile", s.ctx, &dto.ProfileInputDTO{Name: &s.newItem.Name, Permissions: map[string]any{"profile": true}}).Return(&s.newItem, nil)
	repo.On("CreateProfile", s.ctx, nilDTO).Return(nil, errors.New("error to create item"))

	repo.On("UpdateProfile", s.ctx, &s.newItem, &dto.ProfileInputDTO{Name: &s.newItem.Name, Permissions: map[string]any{"profile": true}}).Return(nil)
	invalidName := "."
	repo.On("UpdateProfile", s.ctx, &s.newItem, &dto.ProfileInputDTO{Name: &invalidName, Permissions: map[string]any{"profile": true}}).Return(ErrInvalidValue)

	repo.On("DeleteProfiles", s.ctx, []uint{1}).Return(nil)
	repo.On("DeleteProfiles", s.ctx, []uint{7}).Return(ErrItemNotFound)

	s.service = NewProfileService(repo)
}

func (s *ProfileTestSuite) TestGetProfiles() {
	items, err := s.service.GetProfiles(s.ctx, s.f)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.ProfileOutputDTO]{}, items)
	s.Equal(items.Count, int64(len(s.items)))

	items, err = s.service.GetProfiles(s.ctx, nil)

	s.Error(err)
	s.Nil(items)
}

func (s *ProfileTestSuite) TestGetProfileByID() {
	item, err := s.service.GetProfileByID(s.ctx, 1)

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(*item.Name, s.firstItem.Name)
	s.Equal(*item.ID, s.firstItem.ID)

	item, err = s.service.GetProfileByID(s.ctx, 7)

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

func (s *ProfileTestSuite) TestCreateProfile() {
	data := &dto.ProfileInputDTO{Name: &s.newItem.Name, Permissions: map[string]any{"profile": true}}
	item, err := s.service.CreateProfile(s.ctx, data)

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(*item.Name, s.newItem.Name)
	s.Equal(*item.ID, s.newItem.ID)

	item, err = s.service.CreateProfile(s.ctx, nil)

	s.Error(err)
	s.Equal("error to create item", err.Error())
	s.Nil(item)
}

func (s *ProfileTestSuite) TestUpdateProfile() {
	data := &dto.ProfileInputDTO{Name: &s.newItem.Name, Permissions: map[string]any{"profile": true}}
	item, err := s.service.UpdateProfile(s.ctx, 4, data)

	s.NoError(err)
	s.IsType(&dto.ProfileOutputDTO{}, item)
	s.Equal(*item.Name, s.newItem.Name)
	s.Equal(*item.ID, s.newItem.ID)

	item, err = s.service.UpdateProfile(s.ctx, 7, data)

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)

	invalidName := "."
	data = &dto.ProfileInputDTO{Name: &invalidName, Permissions: map[string]any{"profile": true}}
	item, err = s.service.UpdateProfile(s.ctx, 4, data)

	s.Error(err)
	s.True(errors.Is(err, ErrInvalidValue))
	s.Nil(item)
}

func (s *ProfileTestSuite) TestDeleteProfiles() {
	s.NoError(s.service.DeleteProfiles(s.ctx, []uint{1}))

	err := s.service.DeleteProfiles(s.ctx, []uint{7})

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
}
