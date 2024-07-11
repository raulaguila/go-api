package handler

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/internal/pkg/i18n"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/validator"
)

var middlewareUserDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.UserInputDTO{},
})

var middlewarePasswordDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.PasswordInputDTO{},
})

type UserHandler struct {
	userService domain.UserService
}

func (h *UserHandler) foreignKeyViolatedFrom(c *fiber.Ctx, messages *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrProfileNotFound)
	case fiber.MethodDelete:
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUserUsed)
	default:
		return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
	}
}

func (h *UserHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(helper.LocalLang).(*i18n.Translation)

	switch pgErr := pgutils.HandlerError(err); {
	case errors.Is(pgErr, pgutils.ErrDuplicatedKey):
		return helper.NewHTTPResponse(c, fiber.StatusConflict, messages.ErrUserRegistered)
	case errors.Is(pgErr, pgutils.ErrForeignKeyViolated):
		return h.foreignKeyViolatedFrom(c, messages)
	case errors.Is(pgErr, pgutils.ErrUndefinedColumn):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUndefinedColumn)
	case errors.Is(pgErr, myerrors.ErrUserHasNoPhoto):
		return helper.NewHTTPResponse(c, fiber.StatusNotFound, messages.ErrUserHasNoPhoto)
	case errors.Is(err, myerrors.ErrUserHasPass):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUserHasPass)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return helper.NewHTTPResponse(c, fiber.StatusNotFound, messages.ErrUserNotFound)
	case errors.As(err, &validator.ErrValidator):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

// NewUserHandler Creates a new user handler.
func NewUserHandler(route fiber.Router, us domain.UserService) {
	handler := &UserHandler{
		userService: us,
	}
	extensions := []string{".jpg", ".jpeg", ".png"}

	route.Put("/pass", middlewarePasswordDTO, handler.setUserPassword)
	route.Delete("/pass", middlewareIDDTO, handler.resetUserPassword)

	route.Use(middleware.MidAccess)

	midFiles := middleware.GetFileFromRequest("photo", &extensions)

	route.Get("", middlewareUserFilterDTO, handler.getUsers)
	route.Post("", middlewareUserDTO, handler.createUser)
	route.Get("/:"+helper.ParamID, middlewareIDDTO, handler.getUser)
	route.Put("/:"+helper.ParamID, middlewareIDDTO, middlewareUserDTO, handler.updateUser)
	route.Delete("", handler.deleteUser)
	route.Put("/:"+helper.ParamID+"/photo", middlewareIDDTO, midFiles, handler.setUserPhoto)
	route.Get("/:"+helper.ParamID+"/photo", middlewareIDDTO, handler.getUserPhoto)
}

// getUserPhoto godoc
// @Summary      Get user photo
// @Description  Get user photo
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query	string				false	"Response language" enums(en-US,pt-BR)
// @Param        id				path	filters.IDFilter	true	"User ID"
// @Success      200  {object}  nil
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user/{id}/photo [get]
// @Security	 Bearer
func (h *UserHandler) getUserPhoto(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	url, err := h.userService.GenerateUserPhotoURL(c.Context(), id.ID)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Redirect(url)
}

// setUserPhoto godoc
// @Summary      Set user photo
// @Description  Set user photo
// @Tags         User
// @Accept       mpfd
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        id				path    	filters.IDFilter	true	"User ID"
// @Param		 photo			formData	file				true	"profile photo"
// @Success      200  {object}  nil
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user/{id}/photo [put]
// @Security	 Bearer
func (h *UserHandler) setUserPhoto(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	photo := c.Locals(helper.LocalDTO).(*domain.File)
	if err := h.userService.SetUserPhoto(c.Context(), id.ID, photo); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// getUsers godoc
// @Summary      Get users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        filter			query		filters.UserFilter	false	"Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO[dto.UserOutputDTO]
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	response, err := h.userService.GetUsers(c.Context(), c.Locals(helper.LocalFilter).(*filters.UserFilter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createUser godoc
// @Summary      Insert user
// @Description  Insert user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        user			body		dto.UserInputDTO	true	"User model"
// @Success      201  {object}  dto.UserOutputDTO
// @Failure      400  {object}  helper.HTTPResponse
// @Failure      409  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(helper.LocalDTO).(*dto.UserInputDTO)
	user, err := h.userService.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// getUser godoc
// @Summary      Get user
// @Description  Get user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        id				path		filters.IDFilter	true	"User ID"
// @Success      200  {object}  dto.UserOutputDTO
// @Failure      400  {object}  helper.HTTPResponse
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user/{id} [get]
// @Security	 Bearer
func (h *UserHandler) getUser(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	user, err := h.userService.GetUserByID(c.Context(), id.ID)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// updateUser godoc
// @Summary      Update user
// @Description  Update user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        id				path		filters.IDFilter	true	"User ID"
// @Param        user			body		dto.UserInputDTO	true	"User model"
// @Success      200  {object}  dto.UserOutputDTO
// @Failure      400  {object}  helper.HTTPResponse
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	userDTO := c.Locals(helper.LocalDTO).(*dto.UserInputDTO)
	user, err := h.userService.UpdateUser(c.Context(), id.ID, userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// deleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        id				body		dto.IDsInputDTO		true	"User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return h.handlerError(c, err)
	}

	if err := h.userService.DeleteUsers(c.Context(), toDelete.IDs); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// resetUser godoc
// @Summary      Reset user password
// @Description  Reset user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string				false	"Response language" enums(en-US,pt-BR)
// @Param        email			query		string				true 	"User email"
// @Success      200  {object}  nil
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user/pass [delete]
// @Security	 Bearer
func (h *UserHandler) resetUserPassword(c *fiber.Ctx) error {
	mail := strings.ReplaceAll(c.Query(helper.ParamMail), "%40", "@")
	if err := h.userService.ResetUserPassword(c.Context(), mail); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// passwordUser godoc
// @Summary      Set user password
// @Description  Set user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang			query		string					false	"Response language" enums(en-US,pt-BR)
// @Param        email			query		string					true	"User email" Format(email)
// @Param        password		body		dto.PasswordInputDTO	true	"Password model"
// @Success      200  {object}  nil
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /user/pass [put]
func (h *UserHandler) setUserPassword(c *fiber.Ctx) error {
	pass := c.Locals(helper.LocalDTO).(*dto.PasswordInputDTO)
	if pass.Password == nil || *pass.Password != *pass.PasswordConfirm {
		messages := c.Locals(helper.LocalLang).(*i18n.Translation)
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrPassUnmatch)
	}

	mail := strings.ReplaceAll(c.Query(helper.ParamMail), "%40", "@")
	if err := h.userService.SetUserPassword(c.Context(), mail, pass); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
