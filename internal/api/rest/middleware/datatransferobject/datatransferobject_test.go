package datatransferobject

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/raulaguila/go-api/pkg/packhub"
)

type TestModel struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		body         io.Reader
		query        *string
		config       []Config
		expectedCode int
	}{
		{
			name: "valid body parser",
			body: strings.NewReader(`{"name":"John","age":20}`),
			config: []Config{{
				OnLookup: Body,
				Model:    &TestModel{},
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					return c.SendStatus(fiber.StatusBadRequest)
				},
			}},
			expectedCode: fiber.StatusOK,
		},
		{
			name: "invalid body parser",
			body: strings.NewReader(`{"name":"John","age":"wrong age type"}`),
			config: []Config{{
				OnLookup: Body,
				Model:    &TestModel{},
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					return c.SendStatus(fiber.StatusBadRequest)
				},
			}},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:  "valid query parser",
			query: packhub.Pointer("name=John&age=20"),
			config: []Config{{
				OnLookup: Query,
				Model:    &TestModel{},
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					return c.SendStatus(fiber.StatusBadRequest)
				},
			}},
			expectedCode: fiber.StatusOK,
		},
		{
			name: "default config",
			body: strings.NewReader(`{"name":"John","age":20}`),
			config: []Config{{
				Model: &TestModel{},
			}},
			expectedCode: fiber.StatusOK,
		},
		{
			name: "default config with error",
			body: strings.NewReader(`{"name":"John","age":"wrong age type"}`),
			config: []Config{{
				Model: &TestModel{},
			}},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(New(tt.config...))

			app.Post("/test", func(c *fiber.Ctx) error {
				obj, ok := c.Locals("localDTO").(*TestModel)
				if !ok {
					return c.SendStatus(fiber.StatusBadRequest)
				}
				return c.Status(fiber.StatusOK).JSON(obj)
			})

			endpoint := "/test"
			if tt.query != nil {
				endpoint += "?" + *tt.query
			}
			req := httptest.NewRequest(fiber.MethodPost, endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			resp, err := app.Test(req)
			require.NoError(t, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(t, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}
