package handler

import (
	"github.com/gofiber/fiber/v2"
	datatransferobject2 "github.com/raulaguila/go-api/internal/api/rest/middleware/datatransferobject"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/utils"
)

var middlewareFilterDTO = datatransferobject2.New(datatransferobject2.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject2.Query,
	Model:      &filter.Filter{},
})

var middlewareUserFilterDTO = datatransferobject2.New(datatransferobject2.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject2.Query,
	Model:      &filters.UserFilter{},
})

var middlewareIDDTO = datatransferobject2.New(datatransferobject2.Config{
	ContextKey: utils.LocalID,
	OnLookup:   datatransferobject2.Params,
	Model:      &filters.IDFilter{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return utils.NewHTTPResponse(c, fiber.StatusBadRequest, "invalidID")
	},
})

var middlewareIDsDTO = datatransferobject2.New(datatransferobject2.Config{
	ContextKey: utils.LocalID,
	OnLookup:   datatransferobject2.Body,
	Model:      &dto.IDsInputDTO{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return utils.NewHTTPResponse(c, fiber.StatusBadRequest, "invalidID")
	},
})
