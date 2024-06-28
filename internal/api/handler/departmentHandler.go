package handler

import (
	"errors"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"log"

	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/helper"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/i18n"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/validator"
)

var middlewareDepartmentDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.DepartmentInputDTO{},
})

type DepartmentHandler struct {
	departmentService domain.DepartmentService
}

func (h *DepartmentHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrDepartmentNotFound)
	case fiber.MethodDelete:
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrDepartmentUsed)
	default:
		return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (h *DepartmentHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(helper.LocalLang).(*i18n.Translation)

	switch pgErr := pgutils.HandlerError(err); {
	case errors.Is(pgErr, pgutils.ErrDuplicatedKey):
		return helper.NewHTTPResponse(c, fiber.StatusConflict, translation.ErrDepartmentRegistered)
	case errors.Is(pgErr, pgutils.ErrForeignKeyViolated):
		return h.foreignKeyViolatedMethod(c, translation)
	case errors.Is(pgErr, pgutils.ErrUndefinedColumn):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrUndefinedColumn)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return helper.NewHTTPResponse(c, fiber.StatusNotFound, translation.ErrDepartmentNotFound)
	case errors.As(err, &validator.ErrValidator):
		return helper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
}

// NewDepartmentHandler Creates a new department handler.
func NewDepartmentHandler(route fiber.Router, ps domain.DepartmentService) {
	handler := &DepartmentHandler{
		departmentService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middlewareFilterDTO, handler.getDepartments)
	route.Post("", middlewareDepartmentDTO, handler.createDepartment)
	route.Get("/:"+helper.ParamID, middlewareIDDTO, handler.getDepartmentByID)
	route.Put("/:"+helper.ParamID, middlewareIDDTO, middlewareDepartmentDTO, handler.updateDepartment)
	route.Delete("", handler.deleteDepartments)
}

// getDepartments godoc
// @Summary      List departments
// @Description  List departments
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO[dto.DepartmentOutputDTO]
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /department [get]
// @Security	 Bearer
func (h *DepartmentHandler) getDepartments(c *fiber.Ctx) error {
	response, err := h.departmentService.GetDepartments(c.Context(), c.Locals(helper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getDepartmentByID godoc
// @Summary      Get department by ID
// @Description  Get department by ID
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Department ID"
// @Success      200  {object}  dto.DepartmentOutputDTO
// @Failure      400  {object}  helper.HTTPResponse
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /department/{id} [get]
// @Security	 Bearer
func (h *DepartmentHandler) getDepartmentByID(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	department, err := h.departmentService.GetDepartmentByID(c.Context(), id.ID)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(department)
}

// createDepartment godoc
// @Summary      Insert department
// @Description  Insert department
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        department body dto.DepartmentInputDTO true "Department model"
// @Success      201  {object}  dto.DepartmentOutputDTO
// @Failure      400  {object}  helper.HTTPResponse
// @Failure      409  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /department [post]
// @Security	 Bearer
func (h *DepartmentHandler) createDepartment(c *fiber.Ctx) error {
	departmentDTO := c.Locals(helper.LocalDTO).(*dto.DepartmentInputDTO)
	department, err := h.departmentService.CreateDepartment(c.Context(), departmentDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(department)
}

// updateDepartment godoc
// @Summary      Update department by ID
// @Description  Update department by ID
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Department ID"
// @Param        department body dto.DepartmentInputDTO true "Department model"
// @Success      200  {object}  dto.DepartmentOutputDTO
// @Failure      400  {object}  helper.HTTPResponse
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /department/{id} [put]
// @Security	 Bearer
func (h *DepartmentHandler) updateDepartment(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	departmentDTO := c.Locals(helper.LocalDTO).(*dto.DepartmentInputDTO)
	department, err := h.departmentService.UpdateDepartment(c.Context(), id.ID, departmentDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(department)
}

// deleteDepartments godoc
// @Summary      Delete departments by IDs
// @Description  Delete departments by IDs
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   body      dto.IDsInputDTO     true        "Department ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  helper.HTTPResponse
// @Failure      500  {object}  helper.HTTPResponse
// @Router       /department [delete]
// @Security	 Bearer
func (h *DepartmentHandler) deleteDepartments(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return h.handlerError(c, err)
	}

	if err := h.departmentService.DeleteDepartments(c.Context(), toDelete.IDs); err != nil {
		return h.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
