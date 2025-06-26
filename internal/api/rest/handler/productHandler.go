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
	"github.com/raulaguila/go-api/pkg/pgfilter"
	"github.com/raulaguila/go-api/pkg/utils"
)

var middlewareProductDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalDTO,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.ProductInputDTO{},
})

type ProductHandler struct {
	productService domain.ProductService
	handlerError   func(*fiber.Ctx, error) error
}

func NewProductHandler(route fiber.Router, ps domain.ProductService) {
	handler := &ProductHandler{
		productService: ps,
		handlerError: newErrorHandler(map[string]map[error][]any{
			fiber.MethodDelete: {
				pgerror.ErrForeignKeyViolated: []any{fiber.StatusBadRequest, "productUsed"},
			},
			"*": {
				utils.ErrInvalidID:         []any{fiber.StatusBadRequest, "invalidID"},
				pgerror.ErrUndefinedColumn: []any{fiber.StatusBadRequest, "undefinedColumn"},
				pgerror.ErrDuplicatedKey:   []any{fiber.StatusConflict, "productRegistered"},
				gorm.ErrRecordNotFound:     []any{fiber.StatusNotFound, "productNotFound"},
			},
		}),
	}

	route.Use(middleware.MidAccess)

	route.Get("", middlewareFilterDTO, handler.getProducts)
	route.Post("", middlewareProductDTO, handler.createProduct)
	route.Put("/:"+utils.ParamID, middlewareIDDTO, middlewareProductDTO, handler.updateProduct)
	route.Delete("", middlewareIDsIntDTO, handler.deleteProducts)
}

// getProducts godoc
// @Summary      List products
// @Description  List products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool			false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string			false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        pgfilter				query	pgfilter.Filter	false	"Optional Filter"
// @Success      200  {array}   	dto.ItemsOutputDTO[dto.ProductOutputDTO]
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /product [get]
// @Security	 Bearer
func (s *ProductHandler) getProducts(c *fiber.Ctx) error {
	response, err := s.productService.GetProducts(c.Context(), c.Locals(utils.LocalFilter).(*pgfilter.Filter))
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createProduct godoc
// @Summary      Insert product
// @Description  Insert product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        product 			body	dto.ProductInputDTO		true	"Product model"
// @Success      201  {object}  	dto.ProductOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      409  {object} 		HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /product [post]
// @Security	 Bearer
func (s *ProductHandler) createProduct(c *fiber.Ctx) error {
	productDTO, err := s.productService.CreateProduct(c.Context(), c.Locals(utils.LocalDTO).(*dto.ProductInputDTO))
	if err != nil {
		return s.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusCreated, fiberi18n.MustLocalize(c, "productCreated"), productDTO)
}

// updateProduct godoc
// @Summary      Update product by ID
// @Description  Update product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool				false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string				false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        id					path    dto.IDFilter		true	"Product ID"
// @Param        product 		body dto.ProductInputDTO true "Product model"
// @Success      200  {object}  	dto.ProductOutputDTO
// @Failure      400  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /product/{id} [put]
// @Security	 Bearer
func (s *ProductHandler) updateProduct(c *fiber.Ctx) error {
	id := c.Locals(utils.LocalID).(*dto.IDFilter)
	productDTO, err := s.productService.UpdateProduct(c.Context(), id.ID, c.Locals(utils.LocalDTO).(*dto.ProductInputDTO))
	if err != nil {
		return s.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "productUpdated"), productDTO)
}

// deleteProducts godoc
// @Summary      Delete products by IDs
// @Description  Delete products by IDs
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        X-Skip-Auth		header	bool					false	"Skip auth" enums(true,false) default(true)
// @Param        Accept-Language	header	string					false	"Request language" enums(en-US,pt-BR) default(en-US)
// @Param        ids				body	dto.IDsInputDTO[uint]	true	"Products ID"
// @Success      204  {object}  	HTTPResponse.Response
// @Failure      404  {object}  	HTTPResponse.Response
// @Failure      500  {object}  	HTTPResponse.Response
// @Router       /product [delete]
// @Security	 Bearer
func (s *ProductHandler) deleteProducts(c *fiber.Ctx) error {
	toDelete := c.Locals(utils.LocalID).(*dto.IDsInputDTO[uint])
	if err := s.productService.DeleteProducts(c.Context(), toDelete.IDs); err != nil {
		return s.handlerError(c, err)
	}

	return HTTPResponse.New(c, fiber.StatusOK, fiberi18n.MustLocalize(c, "productDeleted"), nil)
}
