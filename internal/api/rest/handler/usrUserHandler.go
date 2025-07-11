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

type userHandler struct {
	service      domain.UserService
	handlerError func(*fiber.Ctx, error) error
}

func NewUserHandler(route fiber.Router, service domain.UserService) {
	handler := &userHandler{
		service: service,
		handlerError: newErrorHandler(map[string]map[error][]any{
			fiber.MethodDelete: {
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "userUsed"},
			},
			"*": {
				utils.ErrInvalidID:            []any{fiber.StatusBadRequest, "invalidID"},
				utils.ErrUserHasPass:          []any{fiber.StatusBadRequest, "hasPass"},
				utils.ErrPasswordsDoNotMatch:  []any{fiber.StatusBadRequest, "passNotMatch"},
				pgerror.ErrUndefinedColumn:    []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:      []any{fiber.StatusConflict, "userRegistered"},
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusNotFound, "itemNotFound"},
				gorm.ErrRecordNotFound:        []any{fiber.StatusNotFound, "userNotFound"},
			},
		}),
	}

	route.Put("/pass", middlewarePasswordDTO, handler.setUserPassword)

	route.Use(middleware.MidAccess)

	route.Delete("/pass", handler.resetUserPassword)
	route.Get("", middlewareUserFilterDTO, handler.getUsers)
	route.Post("", middlewareUserDTO, handler.createUser)
	route.Put("/:"+utils.ParamID, middlewareIDIntDTO, middlewareUserDTO, handler.updateUser)
	route.Delete("", middlewareIDsIntDTO, handler.deleteUser)
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
func (h *userHandler) getUsers(c *fiber.Ctx) error {
	response, err := h.service.GetUsers(c.Context(), c.Locals(utils.LocalFilter).(*dto.UserFilter))
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
func (h *userHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(utils.LocalDTO).(*dto.UserInputDTO)
	user, err := h.service.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusCreated, fiberi18n.MustLocalize(c, "userCreated"), user)
}

// updateUser godoc
// @Summary      Update user by ID
// @Description  Update user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path		dto.IDFilter[uint]	true	"User ID"
// @Param        user				body		dto.UserInputDTO	true	"User model"
// @Success      200  {object}  	dto.UserOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *userHandler) updateUser(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*dto.IDFilter[uint])
	user, err := h.service.UpdateUser(c.Context(), id.ID, c.Locals(utils.LocalDTO).(*dto.UserInputDTO))
	if err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "userUpdated"), user)
}

// deleteUser godoc
// @Summary      Delete user by ID
// @Description  Delete user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header		bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header		string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					body		dto.IDsInputDTO[uint]	true	"User ID"
// @Success      204  {object}  	nil
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /user [delete]
// @Security	 Bearer
func (h *userHandler) deleteUser(c *fiber.Ctx) error {
	toDelete := c.Locals(utils.LocalID).(*dto.IDsInputDTO[uint])
	if err := h.service.DeleteUsers(c.Context(), toDelete.IDs); err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "userDeleted"), nil)
}

// resetUser godoc
// @Summary      Reset user password by ID
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
func (h *userHandler) resetUserPassword(c *fiber.Ctx) error {
	email, err := url.QueryUnescape(c.Query(utils.ParamMail, ""))
	if err != nil {
		return h.handlerError(c, err)
	}

	if err := h.service.ResetUserPassword(c.Context(), email); err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passReset"), nil)
}

// passwordUser godoc
// @Summary      Set user password by ID
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
func (h *userHandler) setUserPassword(c *fiber.Ctx) error {
	email, err := url.QueryUnescape(c.Query(utils.ParamMail, ""))
	if err != nil {
		return h.handlerError(c, err)
	}

	pass := c.Locals(utils.LocalDTO).(*dto.PasswordInputDTO)
	if pass.Password == nil || pass.PasswordConfirm == nil || *pass.Password != *pass.PasswordConfirm {
		return h.handlerError(c, utils.ErrPasswordsDoNotMatch)
	}

	if err := h.service.SetUserPassword(c.Context(), email, pass); err != nil {
		return h.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "passSet"), nil)
}
