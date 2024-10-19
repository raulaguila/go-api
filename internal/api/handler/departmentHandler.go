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

var middlewareDepartmentDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.DepartmentInputDTO{},
})

type DepartmentHandler struct {
	departmentService domain.DepartmentService
	handlerError      func(*fiber.Ctx, error) error
}

// NewDepartmentHandler Creates a new department handler.
func NewDepartmentHandler(route fiber.Router, ps domain.DepartmentService) {
	localErrors := map[string]map[error][]any{
		"*": {
			pgutils.ErrUndefinedColumn: []any{fiber.StatusBadRequest, "undefinedColumn"},
			pgutils.ErrDuplicatedKey:   []any{fiber.StatusConflict, "departmentRegistered"},
			gorm.ErrRecordNotFound:     []any{fiber.StatusNotFound, "departmentNotFound"},
		},
		fiber.MethodDelete: {
			pgutils.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "departmentUsed"},
		},
	}

	handler := &DepartmentHandler{
		departmentService: ps,
		handlerError:      NewErrorHandler(localErrors),
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
// @Param        Accept-Language	header	string			false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        filter				query	filter.Filter	false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.DepartmentOutputDTO]
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /department [get]
// @Security	 Bearer
func (s *DepartmentHandler) getDepartments(c *fiber.Ctx) error {
	response, err := s.departmentService.GetDepartments(c.Context(), c.Locals(helper.LocalFilter).(*filter.Filter))
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getDepartmentByID godoc
// @Summary      Get department by ID
// @Description  Get department by ID
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path	filters.IDFilter	true	"Department ID"
// @Success      200  {object}  	dto.DepartmentOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /department/{id} [get]
// @Security	 Bearer
func (s *DepartmentHandler) getDepartmentByID(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	department, err := s.departmentService.GetDepartmentByID(c.Context(), id.ID)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(department)
}

// createDepartment godoc
// @Summary      Insert department
// @Description  Insert department
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        department 		body	dto.DepartmentInputDTO	true	"Department model"
// @Success      201  {object}  	dto.DepartmentOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      409  {object} 		helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /department [post]
// @Security	 Bearer
func (s *DepartmentHandler) createDepartment(c *fiber.Ctx) error {
	departmentDTO := c.Locals(helper.LocalDTO).(*dto.DepartmentInputDTO)
	department, err := s.departmentService.CreateDepartment(c.Context(), departmentDTO)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(department)
}

// updateDepartment godoc
// @Summary      Update department by ID
// @Description  Update department by ID
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    filters.IDFilter	true	"Department ID"
// @Param        department 		body dto.DepartmentInputDTO true "Department model"
// @Success      200  {object}  	dto.DepartmentOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /department/{id} [put]
// @Security	 Bearer
func (s *DepartmentHandler) updateDepartment(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	departmentDTO := c.Locals(helper.LocalDTO).(*dto.DepartmentInputDTO)
	department, err := s.departmentService.UpdateDepartment(c.Context(), id.ID, departmentDTO)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(department)
}

// deleteDepartments godoc
// @Summary      Delete departments by IDs
// @Description  Delete departments by IDs
// @Tags         Department
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					body	dto.IDsInputDTO		true	"Department ID"
// @Success      204  {object}  	nil
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /department [delete]
// @Security	 Bearer
func (s *DepartmentHandler) deleteDepartments(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return s.handlerError(c, err)
	}

	if err := s.departmentService.DeleteDepartments(c.Context(), toDelete.IDs); err != nil {
		return s.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
