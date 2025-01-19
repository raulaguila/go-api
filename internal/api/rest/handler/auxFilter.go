package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/api/rest/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/utils"
)

var middlewareFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &filter.Filter{},
})

var middlewareUserFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &filters.UserFilter{},
})

var middlewareIDDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalID,
	OnLookup:   datatransferobject.Params,
	Model:      &filters.IDFilter{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return utils.NewHTTPResponse(c, fiber.StatusBadRequest, "invalidID", nil)
	},
})

var middlewareIDsDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalID,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.IDsInputDTO{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return utils.NewHTTPResponse(c, fiber.StatusBadRequest, "invalidID", nil)
	},
})
