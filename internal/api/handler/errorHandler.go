package handler

import (
	"errors"
	"log"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

// newErrorHandler creates a new error handler for a Fiber application, mapping possible errors to HTTP responses.
// The function takes a nested map of method names to error mappings, and returns an error handling function for Fiber.
// It checks if the error matches predefined errors for the given HTTP method and sends appropriate responses.
// If no error matches, it logs the error and responds with a generic internal server error message.
func newErrorHandler(possiblesErrors map[string]map[error][]any) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		for method, mapper := range possiblesErrors {
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
