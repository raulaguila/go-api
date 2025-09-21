package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/api/rest/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/HTTPResponse"
	"github.com/raulaguila/go-api/internal/pkg/consts"
	"github.com/raulaguila/go-api/internal/pkg/dto"
)

// var middlewareFilterDTO = datatransferobject.New(datatransferobject.Config{
// 	ContextKey: consts.LocalFilter,
// 	OnLookup:   datatransferobject.Query,
// 	Model:      &pgfilter.Filter{},
// })

var middlewareProfileFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: consts.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &dto.ProfileFilter{},
})

var middlewareUserFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: consts.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &dto.UserFilter{},
})

var middlewareIDIntDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: consts.LocalID,
	OnLookup:   datatransferobject.Params,
	Model:      &dto.IDFilter[uint]{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return HTTPResponse.New(c, fiber.StatusBadRequest, "invalidID", nil)
	},
})

var middlewareIDsIntDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: consts.LocalID,
	OnLookup:   datatransferobject.Body,
	Model:      &dto.IDsInputDTO[uint]{},
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return HTTPResponse.New(c, fiber.StatusBadRequest, "invalidID", nil)
	},
})

// var middlewareIDsStringDTO = datatransferobject.New(datatransferobject.Config{
// 	ContextKey: consts.LocalID,
// 	OnLookup:   datatransferobject.Body,
// 	Model:      &dto.IDsInputDTO[string]{},
// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
// 		return HTTPResponse.New(c, fiber.StatusBadRequest, "invalidID", nil)
// 	},
// })
