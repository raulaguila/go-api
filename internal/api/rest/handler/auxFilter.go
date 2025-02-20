package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/api/rest/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/HTTPResponse"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgfilter"
	"github.com/raulaguila/go-api/pkg/utils"
)

var middlewareFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &pgfilter.Filter{},
})

var middlewareProfileFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &dto.ProfileFilter{},
})

var middlewareUserFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &dto.UserFilter{},
})

var middlewareIDDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalID,
	OnLookup:   datatransferobject.Params,
	Model:      &dto.IDFilter{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return HTTPResponse.New(c, fiber.StatusBadRequest, "invalidID", nil)
	},
})

var middlewareIDsDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: utils.LocalID,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.IDsInputDTO{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return HTTPResponse.New(c, fiber.StatusBadRequest, "invalidID", nil)
	},
})
