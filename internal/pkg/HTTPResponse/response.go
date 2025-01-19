package HTTPResponse

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
	Object  any    `json:"object,omitempty"`
}

func New(c *fiber.Ctx, status int, message string, object any) error {
	return c.Status(status).JSON(&Response{
		Code:    status,
		Message: message,
		Object:  object,
	})
}
