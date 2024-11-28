package datatransferobject

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// New creates a new Fiber middleware handler with the given configuration options.
// The middleware extracts data from HTTP requests into a specified struct according to the OnLookup setting.
// By default, the data is parsed from the request body unless otherwise specified in the Config.
// The parsed object is stored in the context using the specified ContextKey for further handling in request lifecycle.
// Optional Config parameters include custom error handling and conditions to skip the middleware.
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	parser := func(c *fiber.Ctx, obj any) (any, error) {
		switch cfg.OnLookup {
		case Body:
			return obj, c.BodyParser(obj)
		case Query:
			return obj, c.QueryParser(obj)
		case Params:
			return obj, c.ParamsParser(obj)
		default:
			return obj, c.CookieParser(obj)
		}
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		var err error
		obj := reflect.New(reflect.TypeOf(cfg.Model).Elem()).Interface()
		obj, err = parser(c, obj)
		if err != nil {
			fmt.Printf("Error mid: %v - %v\n", err, reflect.TypeOf(cfg.Model))
			return cfg.ErrorHandler(c, err)
		}

		c.Locals(cfg.ContextKey, obj)
		return c.Next()
	}
}
