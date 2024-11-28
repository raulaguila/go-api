package helper

import (
	"github.com/gofiber/fiber/v2"
)

// HTTPResponse represents a standard structure for HTTP responses containing a status code and a message.
type HTTPResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// NewHTTPResponse sends an HTTP response with the given status code and message in JSON format using Fiber.
func NewHTTPResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(&HTTPResponse{
		Code:    status,
		Message: message,
	})
}
