package handler

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/api/rest/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/pgfilter"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/utils"
)

var middlewareProfileDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.ProfileInputDTO{},
})

type ProfileHandler struct {
	profileService domain.ProfileService
	handlerError   func(*fiber.Ctx, error) error
}

func NewProfileHandler(route fiber.Router, ps domain.ProfileService) {
	handler := &ProfileHandler{
		profileService: ps,
		handlerError: newErrorHandler(map[string]map[error][]any{
			fiber.MethodDelete: {
				pgutils.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "profileUsed"},
			},
			"*": {
				utils.ErrInvalidID:         []any{fiber.StatusBadRequest, "invalidID"},
				pgutils.ErrUndefinedColumn: []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgutils.ErrDuplicatedKey:   []any{fiber.StatusConflict, "profileRegistered"},
				gorm.ErrRecordNotFound:     []any{fiber.StatusNotFound, "profileNotFound"},
			},
		}),
	}

	route.Use(middleware.MidAccess)

	route.Get("", middlewareFilterDTO, handler.getProfiles)
	route.Post("", middlewareProfileDTO, handler.createProfile)
	route.Get("/:"+utils.ParamID, middlewareIDDTO, handler.getProfile)
	route.Put("/:"+utils.ParamID, middlewareIDDTO, middlewareProfileDTO, handler.updateProfile)
	route.Delete("", middlewareIDsDTO, handler.deleteProfiles)
}

// getProfiles godoc
// @Summary      Get profiles
// @Description  Get profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter				query	pgfilter.Filter		false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.ProfileOutputDTO]
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /profile [get]
// @Security	 Bearer
func (s *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	response, err := s.profileService.GetProfiles(c.Context(), c.Locals(utils.LocalFilter).(*pgfilter.Filter))
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getProfile godoc
// @Summary      Get profile by ID
// @Description  Get profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    filters.IDFilter	true	"Profile ID"
// @Success      200  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	utils.HTTPResponse
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /profile/{id} [get]
// @Security	 Bearer
func (s *ProfileHandler) getProfile(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*filters.IDFilter)
	profile, err := s.profileService.GetProfileByID(c.Context(), id.ID)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// createProfile godoc
// @Summary      Insert profile
// @Description  Insert profile
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        profile			body	dto.ProfileInputDTO	true	"Profile model"
// @Success      201  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	utils.HTTPResponse
// @Failure      409  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /profile [post]
// @Security	 Bearer
func (s *ProfileHandler) createProfile(c *fiber.Ctx) error {
	profileDTO, err := s.profileService.CreateProfile(c.Context(), c.Locals(utils.LocalDTO).(*dto.ProfileInputDTO))
	if err != nil {
		return s.handlerError(c, err)
	}

	return utils.NewHTTPResponse(c, fiber.StatusCreated, fiberi18n.MustLocalize(c, "profileCreated"), profileDTO)
}

// updateProfile godoc
// @Summary      Update profile
// @Description  Update profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    filters.IDFilter	true	"Profile ID"
// @Param        profile			body	dto.ProfileInputDTO true	"Profile model"
// @Success      200  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	utils.HTTPResponse
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /profile/{id} [put]
// @Security	 Bearer
func (s *ProfileHandler) updateProfile(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*filters.IDFilter)
	profileDTO, err := s.profileService.UpdateProfile(c.Context(), id.ID, c.Locals(utils.LocalDTO).(*dto.ProfileInputDTO))
	if err != nil {
		return s.handlerError(c, err)
	}

	return utils.NewHTTPResponse(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "profileUpdated"), profileDTO)
}

// deleteProfiles godoc
// @Summary      Delete profiles by IDs
// @Description  Delete profiles by IDs
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        ids				body	dto.IDsInputDTO     true	"Profiles ID"
// @Success      204  {object}  	utils.HTTPResponse
// @Failure      404  {object}  	utils.HTTPResponse
// @Failure      500  {object}  	utils.HTTPResponse
// @Router       /profile [delete]
// @Security	 Bearer
func (s *ProfileHandler) deleteProfiles(c *fiber.Ctx) error {
	toDelete := c.Locals(utils.LocalID).(*dto.IDsInputDTO)
	if err := s.profileService.DeleteProfiles(c.Context(), toDelete.IDs); err != nil {
		return s.handlerError(c, err)
	}

	return utils.NewHTTPResponse(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "profileDeleted"), nil)
}
