package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewMiscHandler(miscRoute fiber.Router) {
	handler := &MiscHandler{}

	miscRoute.Get("", handler.healthCheck).Name("Root")
}

type MiscHandler struct{}

// healthCheck godoc
// @Summary      Ping Pong
// @Description  Ping Pong
// @Tags         Ping
// @Produce      json
// @Success      200  {object}   map[string]string
// @Router       / [get]
func (h *MiscHandler) healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"time": time.Now(),
	})
}
