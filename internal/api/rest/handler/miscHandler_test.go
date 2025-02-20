package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func setupMiscApp() *fiber.App {
	app := fiber.New()
	NewMiscHandler(app.Group(""))

	return app
}

func TestMiscHandler_healthCheck(t *testing.T) {
	tests := []struct {
		name, endpoint string
		expectedCode   int
	}{
		{
			name:         "success",
			endpoint:     "/",
			expectedCode: fiber.StatusOK,
		},
	}

	app := setupMiscApp()
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, tt.endpoint, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}
