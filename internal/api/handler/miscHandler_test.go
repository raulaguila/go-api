package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestMiscHandler_healthCheck(t *testing.T) {
	app := fiber.New()
	NewMiscHandler(app.Group("/health"))

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{"success", "/health", fiber.StatusOK},
		{"invalid", "/user", fiber.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, tt.route, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
