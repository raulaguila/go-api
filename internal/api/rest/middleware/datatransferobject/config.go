package datatransferobject

import (
	"reflect"

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

// defaultConfig provides the default configuration settings for the application, including context key and error handler.
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

// isPointerOfStruct checks if the provided interface is a pointer to a struct using reflection.
// It returns true if the input is a pointer pointing to a struct, otherwise false.
func isPointerOfStruct(i any) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.TypeOf(i).Elem().Kind() == reflect.Struct
}

// configDefault returns a merged configuration by combining the defaultConfig with the provided custom Config options.
// If no custom Config is provided, it returns the defaultConfig.
// Validates and assigns default values for undefined fields in the provided Config.
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return defaultConfig
	}

	cfg := config[0]
	if cfg.ContextKey == "" {
		cfg.ContextKey = defaultConfig.ContextKey
	}
	if cfg.OnLookup > Cookie {
		cfg.OnLookup = defaultConfig.OnLookup
	}
	if cfg.Model == nil {
		cfg.Model = defaultConfig.Model
	}
	if cfg.OnLookup != Body && !isPointerOfStruct(cfg.Model) {
		panic("model to parse params, queries and cookies must be a pointer of struct")
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = defaultConfig.ErrorHandler
	}

	return cfg
}
