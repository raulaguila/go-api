package handler

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
	"github.com/raulaguila/go-api/pkg/helper"
)

type AuthHandler struct {
	authService  domain.AuthService
	handlerError func(*fiber.Ctx, error) error
}

// NewAuthHandler Creates a new authenticator handler.
func NewAuthHandler(route fiber.Router, as domain.AuthService) {
	localErrors := map[string]map[error][]any{
		"*": {
			myerrors.ErrDisabledUser:       []any{fiber.StatusUnauthorized, "disabledUser"},
			myerrors.ErrInvalidCredentials: []any{fiber.StatusUnauthorized, "incorrectCredentials"},
			gorm.ErrRecordNotFound:         []any{fiber.StatusNotFound, "userNotFound"},
		},
	}

	handler := &AuthHandler{
		authService:  as,
		handlerError: NewErrorHandler(localErrors),
	}

	route.Post("", handler.login)
	route.Get("", middleware.MidAccess, handler.me)
	route.Put("", middleware.MidRefresh, handler.refresh)
}

// login godoc
// @Summary      User authentication
// @Description  User authentication
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        credentials		body	dto.AuthInputDTO	true	"Credentials model"
// @Success      200  {object}  	dto.AuthOutputDTO
// @Failure      401  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /auth [post]
func (s *AuthHandler) login(c *fiber.Ctx) error {
	credentials := &dto.AuthInputDTO{}
	if err := c.BodyParser(credentials); err != nil {
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, fiberi18n.MustLocalize(c, "invalidData"))
	}

	authResponse, err := s.authService.Login(c.Context(), credentials, c.IP())
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(authResponse)
}

// me godoc
// @Summary      User authenticated
// @Description  User authenticated
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization		header	string				false	"User token"
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Success      200  {object}  	dto.UserOutputDTO
// @Failure      401  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /auth [get]
// @Security	 Bearer
func (s *AuthHandler) me(c *fiber.Ctx) error {
	user := c.Locals(helper.LocalUser).(*domain.User)
	return c.Status(fiber.StatusOK).JSON(s.authService.Me(user))
}

// refresh godoc
// @Summary      User refresh
// @Description  User refresh
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization		header	string				false	"User token"
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Success      200  {object}  	dto.AuthOutputDTO
// @Failure      401  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /auth [put]
func (s *AuthHandler) refresh(c *fiber.Ctx) error {
	user := c.Locals(helper.LocalUser).(*domain.User)
	return c.Status(fiber.StatusOK).JSON(s.authService.Refresh(user, c.IP()))
}
