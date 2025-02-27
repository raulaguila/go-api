package handler

import (
	"net/url"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/api/rest/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/HTTPResponse"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgerror"
	"github.com/raulaguila/go-api/pkg/utils"
)

var middlewareUserDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.UserInputDTO{},
})

var middlewarePasswordDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.PasswordInputDTO{},
})

type UserHandler struct {
	userService  domain.UserService
	handlerError func(*fiber.Ctx, error) error
}

func NewUserHandler(route fiber.Router, us domain.UserService) {
	handler := &UserHandler{
		userService: us,
		handlerError: newErrorHandler(map[string]map[error][]any{
			fiber.MethodPost: {
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusNotFound, "itemNotFound"},
			},
			fiber.MethodPut: {
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusNotFound, "itemNotFound"},
			},
			fiber.MethodDelete: {
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "userUsed"},
			},
			"*": {
				utils.ErrInvalidID:           []any{fiber.StatusBadRequest, "invalidID"},
				pgerror.ErrUndefinedColumn:   []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:     []any{fiber.StatusConflict, "userRegistered"},
				utils.ErrUserHasPass:         []any{fiber.StatusBadRequest, "hasPass"},
				utils.ErrPasswordsDoNotMatch: []any{fiber.StatusBadRequest, "passNotMatch"},
				gorm.ErrRecordNotFound:       []any{fiber.StatusNotFound, "userNotFound"},
			},
		}),
	}

	route.Put("/pass", middlewarePasswordDTO, handler.setUserPassword)

	route.Use(middleware.MidAccess)

	route.Delete("/pass", middlewareIDDTO, handler.resetUserPassword)
	route.Get("", middlewareUserFilterDTO, handler.getUsers)
	route.Post("", middlewareUserDTO, handler.createUser)
	route.Get("/:"+utils.ParamID, middlewareIDDTO, handler.getUser)
	route.Put("/:"+utils.ParamID, middlewareIDDTO, middlewareUserDTO, handler.updateUser)
	route.Delete("", middlewareIDsDTO, handler.deleteUser)
}

// getUsers godoc
// @Summary      Get users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query		dto.UserFilter		false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.UserOutputDTO]
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	response, err := h.userService.GetUsers(c.Context(), c.Locals(utils.LocalFilter).(*dto.UserFilter))
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
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      409  {object}  	HTTPResponse.Response
// @Failure      500  {object} 		HTTPResponse.Response
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(utils.LocalDTO).(*dto.UserInputDTO)
	user, err := h.userService.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusCreated, fiberi18n.MustLocalize(c, "userCreated"), user)
}

// getUser godoc
// @Summary      Get user
// @Description  Get user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path		dto.IDFilter		true	"User ID"
// @Success      200  {object}  	dto.UserOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user/{id} [get]
// @Security	 Bearer
func (h *UserHandler) getUser(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*dto.IDFilter)
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
// @Param        id					path		dto.IDFilter		true	"User ID"
// @Param        user				body		dto.UserInputDTO	true	"User model"
// @Success      200  {object}  	dto.UserOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*dto.IDFilter)
	userDTO := c.Locals(utils.LocalDTO).(*dto.UserInputDTO)
	user, err := h.userService.UpdateUser(c.Context(), id.ID, userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "userUpdated"), user)
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
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	toDelete := c.Locals(utils.LocalID).(*dto.IDsInputDTO)
	if err := h.userService.DeleteUsers(c.Context(), toDelete.IDs); err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "userDeleted"), nil)
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
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user/pass [delete]
// @Security	 Bearer
func (h *UserHandler) resetUserPassword(c *fiber.Ctx) error {
	email, err := url.QueryUnescape(c.Query(utils.ParamMail, ""))
	if err != nil {
		return h.handlerError(c, err)
	}

	if err := h.userService.ResetUserPassword(c.Context(), email); err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passReset"), nil)
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
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user/pass [put]
func (h *UserHandler) setUserPassword(c *fiber.Ctx) error {
	email, err := url.QueryUnescape(c.Query(utils.ParamMail, ""))
	if err != nil {
		return h.handlerError(c, err)
	}

	pass := c.Locals(utils.LocalDTO).(*dto.PasswordInputDTO)
	if pass.Password == nil || pass.PasswordConfirm == nil || *pass.Password != *pass.PasswordConfirm {
		return h.handlerError(c, utils.ErrPasswordsDoNotMatch)
	}

	if err := h.userService.SetUserPassword(c.Context(), email, pass); err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passSet"), nil)
}
