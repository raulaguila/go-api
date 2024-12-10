package handler

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type generalHandlerTest struct {
	name, method, endpoint string
	body                   io.Reader
	setupMocks             func()
	expectedCode           int
}

func (s *generalHandlerTest) runTest(t *testing.T, app *fiber.App) {
	t.Run(s.name, func(test *testing.T) {
		s.setupMocks()
		req := httptest.NewRequest(s.method, s.endpoint, s.body)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Set("X-Skip-Auth", "true")

		resp, err := app.Test(req)
		require.NoError(test, err, fmt.Sprintf("Error on test '%v'", s.name))
		require.Equal(test, s.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", s.name))
	})
}

func runGeneralHandlerTests(t *testing.T, tests []generalHandlerTest, app *fiber.App) {
	for _, test := range tests {
		test.runTest(t, app)
	}
}
