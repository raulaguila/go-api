package handler

import (
	"errors"
	"log"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

// NewErrorHandler Generic function that receives a map with the http methods, errors, status code and the message for each error.
func NewErrorHandler(possibleErrors map[string]map[error][]any) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		for method, mapper := range possibleErrors {
			if method == c.Method() || method == "*" {
				for key, value := range mapper {
					switch pgErr := pgutils.HandlerError(err); {
					case errors.Is(pgErr, key):
						return helper.NewHTTPResponse(c, value[0].(int), fiberi18n.MustLocalize(c, value[1].(string)))
					}
				}
			}
		}

		log.Printf("Undected error: %s\n", err.Error())
		return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, fiberi18n.MustLocalize(c, "errGeneric"))
	}
}
