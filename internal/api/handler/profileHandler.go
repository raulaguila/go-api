package handler

import (
	"errors"
	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/validator"
	"gorm.io/gorm"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/i18n"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/pg-utils"
)

var middlewareProfileDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.ProfileInputDTO{},
})

type ProfileHandler struct {
	profileService domain.ProfileService
}

func (h *ProfileHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrProfileNotFound)
	case fiber.MethodDelete:
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrProfileUsed)
	default:
		return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (h *ProfileHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(helper.LocalLang).(*i18n.Translation)

	switch pgErr := pg_utils.HandlerError(err); {
	case errors.Is(pgErr, pg_utils.ErrDuplicatedKey):
		return helper.NewHTTPResponse(c, fiber.StatusConflict, messages.ErrProfileRegistered)
	case errors.Is(pgErr, pg_utils.ErrForeignKeyViolated):
		return h.foreignKeyViolatedMethod(c, messages)
	case errors.Is(pgErr, pg_utils.ErrUndefinedColumn):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUndefinedColumn)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return helper.NewHTTPResponse(c, fiber.StatusNotFound, messages.ErrProfileNotFound)
	case errors.As(err, &validator.ErrValidator):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

// NewProfileHandler Creates a new profile handler.
func NewProfileHandler(route fiber.Router, ps domain.ProfileService) {
	handler := &ProfileHandler{
		profileService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middlewareFilterDTO, handler.getProfiles)
	route.Post("", middlewareProfileDTO, handler.createProfile)
	route.Get("/:"+helper.ParamID, middlewareIDDTO, handler.getProfile)
	route.Put("/:"+helper.ParamID, middlewareIDDTO, middlewareProfileDTO, handler.updateProfile)
	route.Delete("", handler.deleteProfiles)
}

// getProfiles godoc
// @Summary      Get profiles
// @Description  Get profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO[dto.ProfileOutputDTO]
// @Failure      500  {object}  http_helper.HTTPResponse
// @Router       /profile [get]
// @Security	 Bearer
func (h *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	response, err := h.profileService.GetProfiles(c.Context(), c.Locals(helper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createProfile godoc
// @Summary      Insert profile
// @Description  Insert profile
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        profile body dto.ProfileInputDTO true "Profile model"
// @Success      201  {object}  dto.ProfileOutputDTO
// @Failure      400  {object}  http_helper.HTTPResponse
// @Failure      409  {object}  http_helper.HTTPResponse
// @Failure      500  {object}  http_helper.HTTPResponse
// @Router       /profile [post]
// @Security	 Bearer
func (h *ProfileHandler) createProfile(c *fiber.Ctx) error {
	profileDTO := c.Locals(helper.LocalDTO).(*dto.ProfileInputDTO)
	profile, err := h.profileService.CreateProfile(c.Context(), profileDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(profile)
}

// getProfile godoc
// @Summary      Get profile by ID
// @Description  Get profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Profile ID"
// @Success      200  {object}  dto.ProfileOutputDTO
// @Failure      400  {object}  http_helper.HTTPResponse
// @Failure      404  {object}  http_helper.HTTPResponse
// @Failure      500  {object}  http_helper.HTTPResponse
// @Router       /profile/{id} [get]
// @Security	 Bearer
func (h *ProfileHandler) getProfile(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	profile, err := h.profileService.GetProfileByID(c.Context(), id.ID)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// updateProfile godoc
// @Summary      Update profile
// @Description  Update profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Profile ID"
// @Param        profile body dto.ProfileInputDTO true "Profile model"
// @Success      200  {object}  dto.ProfileOutputDTO
// @Failure      400  {object}  http_helper.HTTPResponse
// @Failure      404  {object}  http_helper.HTTPResponse
// @Failure      500  {object}  http_helper.HTTPResponse
// @Router       /profile/{id} [put]
// @Security	 Bearer
func (h *ProfileHandler) updateProfile(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	profileDTO := c.Locals(helper.LocalDTO).(*dto.ProfileInputDTO)
	profile, err := h.profileService.UpdateProfile(c.Context(), id.ID, profileDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// deleteProfiles godoc
// @Summary      Delete profiles by IDs
// @Description  Delete profiles by IDs
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   body      dto.IDsInputDTO     true        "Profile ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  http_helper.HTTPResponse
// @Failure      500  {object}  http_helper.HTTPResponse
// @Router       /profile [delete]
// @Security	 Bearer
func (h *ProfileHandler) deleteProfiles(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return h.handlerError(c, err)
	}

	if err := h.profileService.DeleteProfiles(c.Context(), toDelete.IDs); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
