package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

var middlewareProductDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.ProductInputDTO{},
})

type ProductHandler struct {
	productService domain.ProductService
	handlerError   func(*fiber.Ctx, error) error
}

// NewProductHandler Creates a new product handler.
func NewProductHandler(route fiber.Router, ps domain.ProductService) {
	handler := &ProductHandler{
		productService: ps,
		handlerError: newErrorHandler(map[string]map[error][]any{
			fiber.MethodDelete: {
				pgutils.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "productUsed"},
			},
			"*": {
				myerrors.ErrInvalidID:      []any{fiber.StatusBadRequest, "invalidID"},
				pgutils.ErrUndefinedColumn: []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgutils.ErrDuplicatedKey:   []any{fiber.StatusConflict, "productRegistered"},
				gorm.ErrRecordNotFound:     []any{fiber.StatusNotFound, "productNotFound"},
			},
		}),
	}

	//route.Use(middleware.MidAccess)

	route.Get("", middlewareFilterDTO, handler.getProducts)
	route.Post("", middlewareProductDTO, handler.createProduct)
	route.Get("/:"+helper.ParamID, middlewareIDDTO, handler.getProductByID)
	route.Put("/:"+helper.ParamID, middlewareIDDTO, middlewareProductDTO, handler.updateProduct)
	route.Delete("", handler.deleteProducts)
}

// getProducts godoc
// @Summary      List products
// @Description  List products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string			false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        filter				query	filter.Filter	false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.ProductOutputDTO]
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /product [get]
// @Security	 Bearer
func (s *ProductHandler) getProducts(c *fiber.Ctx) error {
	response, err := s.productService.GetProducts(c.Context(), c.Locals(helper.LocalFilter).(*filter.Filter))
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getProductByID godoc
// @Summary      Get product by ID
// @Description  Get product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path	filters.IDFilter	true	"Product ID"
// @Success      200  {object}  	dto.ProductOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /product/{id} [get]
// @Security	 Bearer
func (s *ProductHandler) getProductByID(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	if id.ID == 0 {
		return s.handlerError(c, myerrors.ErrInvalidID)
	}

	product, err := s.productService.GetProductByID(c.Context(), id.ID)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// createProduct godoc
// @Summary      Insert product
// @Description  Insert product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        product 		body	dto.ProductInputDTO	true	"Product model"
// @Success      201  {object}  	dto.ProductOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      409  {object} 		helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /product [post]
// @Security	 Bearer
func (s *ProductHandler) createProduct(c *fiber.Ctx) error {
	productDTO := c.Locals(helper.LocalDTO).(*dto.ProductInputDTO)
	product, err := s.productService.CreateProduct(c.Context(), productDTO)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// updateProduct godoc
// @Summary      Update product by ID
// @Description  Update product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    filters.IDFilter	true	"Product ID"
// @Param        product 		body dto.ProductInputDTO true "Product model"
// @Success      200  {object}  	dto.ProductOutputDTO
// @Failure      400  {object}  	helper.HTTPResponse
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /product/{id} [put]
// @Security	 Bearer
func (s *ProductHandler) updateProduct(c *fiber.Ctx) error {
	id := c.Locals(helper.LocalID).(*filters.IDFilter)
	if id.ID == 0 {
		return s.handlerError(c, myerrors.ErrInvalidID)
	}

	productDTO := c.Locals(helper.LocalDTO).(*dto.ProductInputDTO)
	product, err := s.productService.UpdateProduct(c.Context(), id.ID, productDTO)
	fmt.Printf("[%v] product: %v - %v\n", id.ID, product, err)
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// deleteProducts godoc
// @Summary      Delete products by IDs
// @Description  Delete products by IDs
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					body	dto.IDsInputDTO		true	"Product ID"
// @Success      204  {object}  	nil
// @Failure      404  {object}  	helper.HTTPResponse
// @Failure      500  {object}  	helper.HTTPResponse
// @Router       /product [delete]
// @Security	 Bearer
func (s *ProductHandler) deleteProducts(c *fiber.Ctx) error {
	toDelete := &dto.IDsInputDTO{}
	if err := c.BodyParser(toDelete); err != nil {
		return s.handlerError(c, err)
	}

	if err := s.productService.DeleteProducts(c.Context(), toDelete.IDs); err != nil {
		return s.handlerError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
