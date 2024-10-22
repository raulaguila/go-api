package handler

import (
	"errors"
	"log"

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
	authService domain.AuthService
}

func (s *AuthHandler) handlerError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return helper.NewHTTPResponse(c, fiber.StatusUnauthorized, fiberi18n.MustLocalize(c, "userNotFound"))
	case errors.Is(err, myerrors.ErrInvalidCredentials):
		return helper.NewHTTPResponse(c, fiber.StatusUnauthorized, fiberi18n.MustLocalize(c, "incorrectCredentials"))
	case errors.Is(err, myerrors.ErrDisabledUser):
		return helper.NewHTTPResponse(c, fiber.StatusUnauthorized, fiberi18n.MustLocalize(c, "disabledUser"))
	}

	log.Println(err.Error())
	return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, fiberi18n.MustLocalize(c, "errGeneric"))
}

// NewAuthHandler Creates a new authenticator handler.
func NewAuthHandler(route fiber.Router, as domain.AuthService) {
	handler := &AuthHandler{
		authService: as,
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
