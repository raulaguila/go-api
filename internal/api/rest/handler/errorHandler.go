package handler

import (
	"errors"
	"log"
	"reflect"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/pkg/pgerror"
	"github.com/raulaguila/go-api/pkg/utils"
)

func newErrorHandler(possiblesErrors map[string]map[error][]any) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		for method, mapper := range possiblesErrors {
			if method == c.Method() || method == "*" {
				for key, value := range mapper {
					switch pgErr := pgerror.HandlerError(err); {
					case errors.Is(pgErr, key):
						return utils.NewHTTPResponse(c, value[0].(int), fiberi18n.MustLocalize(c, value[1].(string)), nil)
					}
				}
			}
		}

		log.Printf("Undected error '%v': %s\n", reflect.TypeOf(err), err.Error())
		return utils.NewHTTPResponse(c, fiber.StatusInternalServerError, fiberi18n.MustLocalize(c, "errGeneric"), nil)
	}
}
