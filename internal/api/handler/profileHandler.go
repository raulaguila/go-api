package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-api/internal/api/middleware"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

var middlewareProfileDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.ProfileInputDTO{},
})

type ProfileHandler struct {
	profileService domain.ProfileService
	handlerError   func(*fiber.Ctx, error) error
}

// NewProfileHandler Creates a new profile handler.
func NewProfileHandler(route fiber.Router, ps domain.ProfileService) {
	localErrors := map[string]map[error][]any{
		"*": {
			pgutils.ErrUndefinedColumn: []any{fiber.StatusBadRequest, "undefinedColumn"},
			pgutils.ErrDuplicatedKey:   []any{fiber.StatusConflict, "profileRegistered"},
			gorm.ErrRecordNotFound:     []any{fiber.StatusNotFound, "profileNotFound"},
		},
		fiber.MethodDelete: {
			pgutils.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "profileUsed"},
		},
	}

	handler := &ProfileHandler{
		profileService: ps,
		handlerError:   NewErrorHandler(localErrors),
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
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        filter				query	filter.Filter		false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.ProfileOutputDTO]
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /profile [get]
// @Security	 Bearer
func (s *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	response, err := s.profileService.GetProfiles(c.Context(), c.Locals(helper.LocalFilter).(*filter.Filter))
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createProfile godoc
// @Summary      Insert profile
// @Description  Insert profile
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        profile			body	dto.ProfileInputDTO	true	"Profile model"
// @Success      201  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      409  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /profile [post]
// @Security	 Bearer
func (s *ProfileHandler) createProfile(c *fiber.Ctx) error {
	profileDTO := c.Locals(helper.LocalDTO).(*dto.ProfileInputDTO)
	profile, err := s.profileService.CreateProfile(c.Context(), profileDTO)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(profile)
}

// getProfile godoc
// @Summary      Get profile by ID
// @Description  Get profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    filters.IDFilter	true	"Profile ID"
// @Success      200  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /profile/{id} [get]
// @Security	 Bearer
func (s *ProfileHandler) getProfile(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	profile, err := s.profileService.GetProfileByID(c.Context(), id.ID)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// updateProfile godoc
// @Summary      Update profile
// @Description  Update profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    filters.IDFilter	true	"Profile ID"
// @Param        profile			body	dto.ProfileInputDTO true	"Profile model"
// @Success      200  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /profile/{id} [put]
// @Security	 Bearer
func (s *ProfileHandler) updateProfile(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	profileDTO := c.Locals(helper.LocalDTO).(*dto.ProfileInputDTO)
	profile, err := s.profileService.UpdateProfile(c.Context(), id.ID, profileDTO)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// deleteProfiles godoc
// @Summary      Delete profiles by IDs
// @Description  Delete profiles by IDs
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					body	dto.IDsInputDTO     true	"Profile ID"
// @Success      204  {object}  	nil
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /profile [delete]
// @Security	 Bearer
func (s *ProfileHandler) deleteProfiles(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return s.handlerError(c, err)
	}

	if err := s.profileService.DeleteProfiles(c.Context(), toDelete.IDs); err != nil {
		return s.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
