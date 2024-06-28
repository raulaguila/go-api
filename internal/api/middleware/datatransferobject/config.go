package datatransferobject

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	// ContextKey string key to store the dto object into context.
	// Optional. Default: "localDTO".
	ContextKey string

	// OnLookup Lookup value used to indicate where to extract the request's DTO.
	// Optional. Default value Body.
	// Possible values:
	// - Body
	// - Query
	// - Params
	// - Cookie
	OnLookup Lookup

	// Model pointer to struct to parse dto.
	// Optional. Default value *map[string]any.
	Model any

	// Next defines a function to skip middleware.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// ErrorHandler defines a function which is executed for an invalid key.
	// It may be used to define a custom error.
	// Optional. Default: 400 err.error()
	ErrorHandler fiber.ErrorHandler
}

var defaultConfig = Config{
	ContextKey: "localDTO",
	OnLookup:   Body,
	Model:      new(map[string]any),
	Next:       nil,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	},
}
