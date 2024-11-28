package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GenerateUserPhotoURL(ctx context.Context, id uint) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) SetUserPhoto(ctx context.Context, id uint, file *domain.File) error {
	args := m.Called(ctx, id, file)
	return args.Error(0)
}

func (m *MockUserService) GetUsers(ctx context.Context, filter *filters.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*dto.ItemsOutputDTO[dto.UserOutputDTO]), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, userDTO *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	args := m.Called(ctx, userDTO)
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*dto.UserOutputDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uint, userDTO *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	args := m.Called(ctx, id, userDTO)
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

func (m *MockUserService) DeleteUsers(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockUserService) ResetUserPassword(ctx context.Context, mail string) error {
	args := m.Called(ctx, mail)
	return args.Error(0)
}

func (m *MockUserService) SetUserPassword(ctx context.Context, mail string, pass *dto.PasswordInputDTO) error {
	args := m.Called(ctx, mail, pass)
	return args.Error(0)
}

func TestUserHandler_getUserPhoto(t *testing.T) {
	app := fiber.New()
	mockService := MockUserService{}

	idFilter := &filters.IDFilter{ID: 1}
	ctx := context.Background()
	mockService.On("GenerateUserPhotoURL", ctx, uint(1)).Return("http://example.com/photo.jpg", nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Get("/photo", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalID, idFilter)
		c.SetUserContext(ctx)
		return h.getUserPhoto(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/photo", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "http://example.com/photo.jpg", resp.Header.Get("Location"))
}

func TestUserHandler_setUserPhoto(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	idFilter := &filters.IDFilter{ID: 1}
	file := &domain.File{}
	ctx := context.Background()
	mockService.On("SetUserPhoto", ctx, uint(1), file).Return(nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Put("/photo", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalID, idFilter)
		c.Locals(helper.LocalDTO, file)
		c.SetUserContext(ctx)
		return h.setUserPhoto(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/photo", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_getUsers(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	userFilter := &filters.UserFilter{}
	users := &dto.ItemsOutputDTO[dto.UserOutputDTO]{}
	ctx := context.Background()
	mockService.On("GetUsers", ctx, userFilter).Return(users, nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Get("/users", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalFilter, userFilter)
		c.SetUserContext(ctx)
		return h.getUsers(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_createUser(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	userDTO := &dto.UserInputDTO{}
	user := &dto.UserOutputDTO{}
	ctx := context.Background()
	mockService.On("CreateUser", ctx, userDTO).Return(user, nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Post("/user", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalDTO, userDTO)
		c.SetUserContext(ctx)
		return h.createUser(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUserHandler_getUser(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	idFilter := &filters.IDFilter{ID: 1}
	user := &dto.UserOutputDTO{}
	ctx := context.Background()
	mockService.On("GetUserByID", ctx, uint(1)).Return(user, nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Get("/user", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalID, idFilter)
		c.SetUserContext(ctx)
		return h.getUser(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_updateUser(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	idFilter := &filters.IDFilter{ID: 1}
	userDTO := &dto.UserInputDTO{}
	user := &dto.UserOutputDTO{}
	ctx := context.Background()
	mockService.On("UpdateUser", ctx, uint(1), userDTO).Return(user, nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Put("/user", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalID, idFilter)
		c.Locals(helper.LocalDTO, userDTO)
		c.SetUserContext(ctx)
		return h.updateUser(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/user", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_deleteUser(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	idsDTO := &dto.IDsInputDTO{IDs: []uint{1, 2}}
	ctx := context.Background()
	mockService.On("DeleteUsers", ctx, idsDTO.IDs).Return(nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Delete("/user", func(c *fiber.Ctx) error {
		c.BodyParser(idsDTO)
		c.SetUserContext(ctx)
		return h.deleteUser(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/user", strings.NewReader(`{"IDs":[1,2]}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserHandler_resetUserPassword(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	mail := "user@example.com"
	ctx := context.Background()
	mockService.On("ResetUserPassword", ctx, mail).Return(nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Delete("/password", func(c *fiber.Ctx) error {
		c.Query(helper.ParamMail, url.QueryEscape(mail))
		c.SetUserContext(ctx)
		return h.resetUserPassword(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/password", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_setUserPassword(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	passDTO := &dto.PasswordInputDTO{
		Password:        helper.String("password"),
		PasswordConfirm: helper.String("password"),
	}
	mail := "user@example.com"
	ctx := context.Background()
	mockService.On("SetUserPassword", ctx, mail, passDTO).Return(nil)

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		},
	}

	app.Put("/password", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalDTO, passDTO)
		c.Query(helper.ParamMail, url.QueryEscape(mail))
		c.SetUserContext(ctx)
		return h.setUserPassword(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/password", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_setUserPassword_PasswordsDoNotMatch(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)

	passDTO := &dto.PasswordInputDTO{
		Password:        helper.String("password"),
		PasswordConfirm: helper.String("different"),
	}
	ctx := context.Background()

	h := &UserHandler{
		userService: mockService,
		handlerError: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		},
	}

	app.Put("/password", func(c *fiber.Ctx) error {
		c.Locals(helper.LocalDTO, passDTO)
		c.SetUserContext(ctx)
		return h.setUserPassword(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/password", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
