package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/utils"
)

// middlewareUserDTO is a Fiber middleware configured to parse and store UserInputDTO data from the request body into context.
var middlewareUserDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.UserInputDTO{},
})

// middlewarePasswordDTO is a middleware configuration that extracts and parses password data from HTTP request bodies.
// The parsed data is stored in the request context under a specified key for use in further processing.
var middlewarePasswordDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.PasswordInputDTO{},
})

// UserHandler is responsible for handling user-related HTTP routes and delegating operations to the UserService.
type UserHandler struct {
	userService  domain.UserService
	handlerError func(*fiber.Ctx, error) error
}

// NewUserHandler sets up the routes for user operations and returns a UserHandler instance.
func NewUserHandler(route fiber.Router, us domain.UserService) {
	handler := &UserHandler{
		userService: us,
		handlerError: newErrorHandler(map[string]map[error][]any{
			fiber.MethodPost: {
				pgutils.ErrForeignKeyViolated: []any{fiber.StatusNotFound, "itemNotFound"},
			},
			fiber.MethodPut: {
				pgutils.ErrForeignKeyViolated: []any{fiber.StatusNotFound, "itemNotFound"},
			},
			fiber.MethodDelete: {
				pgutils.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "userUsed"},
			},
			"*": {
				myerrors.ErrInvalidID:           []any{fiber.StatusBadRequest, "invalidID"},
				pgutils.ErrUndefinedColumn:      []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgutils.ErrDuplicatedKey:        []any{fiber.StatusConflict, "userRegistered"},
				myerrors.ErrUserHasPass:         []any{fiber.StatusBadRequest, "hasPass"},
				myerrors.ErrPasswordsDoNotMatch: []any{fiber.StatusBadRequest, "passNotMatch"},
				myerrors.ErrUserHasNoPhoto:      []any{fiber.StatusNotFound, "hasNoPhoto"},
				gorm.ErrRecordNotFound:          []any{fiber.StatusNotFound, "userNotFound"},
			},
		}),
	}

	route.Put("/pass", middlewarePasswordDTO, handler.setUserPassword)
	route.Delete("/pass", middlewareIDDTO, handler.resetUserPassword)

	route.Use(middleware.MidAccess)

	route.Get("", middlewareUserFilterDTO, handler.getUsers)
	route.Post("", middlewareUserDTO, handler.createUser)
	route.Get("/:"+utils.ParamID, middlewareIDDTO, handler.getUser)
	route.Put("/:"+utils.ParamID, middlewareIDDTO, middlewareUserDTO, handler.updateUser)
	route.Delete("", handler.deleteUser)
}

// getUsers godoc
// @Summary      Get users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        filter				query		filters.UserFilter	false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.UserOutputDTO]
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	response, err := h.userService.GetUsers(c.Context(), c.Locals(utils.LocalFilter).(*filters.UserFilter))
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
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        user				body		dto.UserInputDTO	true	"User model"
// @Success      201  {object}  	dto.UserOutputDTO
// @Failure      400  {object}  	utils.HTTPResponse
// @Failure      409  {object}  	utils.HTTPResponse
// @Failure      500  {object} 		utils.HTTPResponse
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(utils.LocalDTO).(*dto.UserInputDTO)
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
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path		filters.IDFilter	true	"User ID"
// @Success      200  {object}  	dto.UserOutputDTO
// @Failure      400  {object}  	utils.HTTPResponse
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /user/{id} [get]
// @Security	 Bearer
func (h *UserHandler) getUser(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*filters.IDFilter)
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
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path		filters.IDFilter	true	"User ID"
// @Param        user				body		dto.UserInputDTO	true	"User model"
// @Success      200  {object}  	dto.UserOutputDTO
// @Failure      400  {object}  	utils.HTTPResponse
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*filters.IDFilter)
	userDTO := c.Locals(utils.LocalDTO).(*dto.UserInputDTO)
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
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					body		dto.IDsInputDTO		true	"User ID"
// @Success      204  {object}  	nil
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /user [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return h.handlerError(c, myerrors.ErrInvalidID)
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
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        email				query		string				true 	"User email"
// @Success      200  {object}  	nil
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /user/pass [delete]
// @Security	 Bearer
func (h *UserHandler) resetUserPassword(c *fiber.Ctx) error {
	mail := strings.ReplaceAll(c.Query(utils.ParamMail), "%40", "@")
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
// @Param        X-Skip-Auth		header		bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        email				query		string					true	"User email" Format(email)
// @Param        password			body		dto.PasswordInputDTO	true	"Password model"
// @Success      200  {object}  	nil
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /user/pass [put]
func (h *UserHandler) setUserPassword(c *fiber.Ctx) error {
	pass := c.Locals(utils.LocalDTO).(*dto.PasswordInputDTO)
	if pass.Password == nil || pass.PasswordConfirm == nil || *pass.Password != *pass.PasswordConfirm {
		return h.handlerError(c, myerrors.ErrPasswordsDoNotMatch)
	}

	mail := strings.ReplaceAll(c.Query(utils.ParamMail), "%40", "@")
	if err := h.userService.SetUserPassword(c.Context(), mail, pass); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
