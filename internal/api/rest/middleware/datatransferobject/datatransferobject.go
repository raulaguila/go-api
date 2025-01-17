package datatransferobject

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

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
