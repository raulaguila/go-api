package datatransferobject

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func isPointerOfStruct(i any) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.TypeOf(i).Elem().Kind() == reflect.Struct
}

func New(config ...Config) fiber.Handler {
	cfg := defaultConfig
	if len(config) > 0 {
		aux := config[0]
		if aux.ContextKey != "" {
			cfg.ContextKey = aux.ContextKey
		}
		if aux.OnLookup <= Cookie {
			cfg.OnLookup = aux.OnLookup
		}
		if aux.Model != nil {
			cfg.Model = aux.Model
			if cfg.OnLookup != Body && !isPointerOfStruct(aux.Model) {
				panic("model to parse params, queries and cookies must be a pointer of struct")
			}
		}
		if aux.Next != nil {
			cfg.Next = aux.Next
		}
		if aux.ErrorHandler != nil {
			cfg.ErrorHandler = aux.ErrorHandler
		}
	}

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
			return cfg.ErrorHandler(c, err)
		}

		c.Locals(cfg.ContextKey, obj)
		return c.Next()
	}
}
