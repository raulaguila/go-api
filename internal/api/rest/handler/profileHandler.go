package handler

import (
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
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "profileUsed"},
			},
			"*": {
				utils.ErrInvalidID:         []any{fiber.StatusBadRequest, "invalidID"},
				pgerror.ErrUndefinedColumn: []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:   []any{fiber.StatusConflict, "profileRegistered"},
				gorm.ErrRecordNotFound:     []any{fiber.StatusNotFound, "profileNotFound"},
			},
		}),
	}

	route.Use(middleware.MidAccess)

	route.Get("", middlewareProfileFilterDTO, handler.getProfiles)
	route.Post("", middlewareProfileDTO, handler.createProfile)
	route.Put("/:"+utils.ParamID, middlewareIDDTO, middlewareProfileDTO, handler.updateProfile)
	route.Delete("", middlewareIDsIntDTO, handler.deleteProfiles)
}

// getProfiles godoc
// @Summary      Get profiles
// @Description  Get profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter			query	dto.ProfileFilter	false	"Profile Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.ProfileOutputDTO]
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /profile [get]
// @Security	 Bearer
func (s *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	f := c.Locals(utils.LocalFilter).(*dto.ProfileFilter)
	f.ListRoot = false
	if u := c.Locals(utils.LocalUser); u != nil && u.(*domain.User).Auth != nil {
		f.ListRoot = u.(*domain.User).Auth.ProfileID == 1
	}

	response, err := s.profileService.GetProfiles(c.Context(), f)
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
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        profile			body	dto.ProfileInputDTO	true	"Profile model"
// @Success      201  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      409  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /profile [post]
// @Security	 Bearer
func (s *ProfileHandler) createProfile(c *fiber.Ctx) error {
	profileDTO, err := s.profileService.CreateProfile(c.Context(), c.Locals(utils.LocalDTO).(*dto.ProfileInputDTO))
	if err != nil {
		return s.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusCreated, fiberi18n.MustLocalize(c, "profileCreated"), profileDTO)
}

// updateProfile godoc
// @Summary      Update profile
// @Description  Update profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    dto.IDFilter		true	"Profile ID"
// @Param        profile			body	dto.ProfileInputDTO true	"Profile model"
// @Success      200  {object}  	dto.ProfileOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /profile/{id} [put]
// @Security	 Bearer
func (s *ProfileHandler) updateProfile(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*dto.IDFilter)
	profileDTO, err := s.profileService.UpdateProfile(c.Context(), id.ID, c.Locals(utils.LocalDTO).(*dto.ProfileInputDTO))
	if err != nil {
		return s.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "profileUpdated"), profileDTO)
}

// deleteProfiles godoc
// @Summary      Delete profiles by IDs
// @Description  Delete profiles by IDs
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        ids				body	dto.IDsInputDTO[uint]   true	"Profiles ID"
// @Success      204  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /profile [delete]
// @Security	 Bearer
func (s *ProfileHandler) deleteProfiles(c *fiber.Ctx) error {
	toDelete := c.Locals(utils.LocalID).(*dto.IDsInputDTO[uint])
	if err := s.profileService.DeleteProfiles(c.Context(), toDelete.IDs); err != nil {
		return s.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "profileDeleted"), nil)
}
