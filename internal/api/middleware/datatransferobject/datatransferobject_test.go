package datatransferobject

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type (
	ProductDTO struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity uint64  `json:"quantity"`
	}

	Product struct {
		ID string `json:"id"`
		ProductDTO
	}

	HttpErrorResponse struct {
		Code    uint64 `json:"code"`
		Message string `json:"message"`
	}
)

const (
	endpoint   = "/product"
	contextKey = "localDTO"
	errorBody  = `{"code":400,"message":"error to decode dto"}`
)

func createAppWithMiddleware2Body() *fiber.App {
	// Create a new app
	app := fiber.New()

	// Add middleware to app
	app.Use(New(Config{
		ContextKey: contextKey,
		OnLookup:   Body,
		Model:      &ProductDTO{},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusBadRequest).JSON(&HttpErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "error to decode dto",
			})
		},
	}))

	// Create a handler to path "/product"
	app.Post(endpoint, func(c *fiber.Ctx) error {
		product := &Product{ID: uuid.New().String(), ProductDTO: *c.Locals(contextKey).(*ProductDTO)}
		return c.Status(fiber.StatusCreated).JSON(product)
	})

	return app
}

// go test -run Test_WithoutBodyDTO
func Test_WithoutBodyDTO(t *testing.T) {
	// Create a new app
	app := createAppWithMiddleware2Body()

	// Create a request without body
	req, _ := http.NewRequest(fiber.MethodPost, endpoint, nil)
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(fiber.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

	// Send request to the app
	res, err := app.Test(req)
	require.NoError(t, err)

	// Read the response body into a string
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	// Check that the response has the expected status code and body
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, errorBody, string(body))
}

// go test -run Test_WithCorrectBodyDTO
func Test_WithCorrectBodyDTO(t *testing.T) {
	// Create a new app
	app := createAppWithMiddleware2Body()

	// Create a request with correct body
	payload := `{"name":"computer","price":5000.5,"quantity":100}`
	req, _ := http.NewRequest(fiber.MethodPost, endpoint, strings.NewReader(payload))
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(fiber.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

	// Send request to the app
	res, err := app.Test(req)
	require.NoError(t, err)

	// Read the response body into a string
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	// Check that the response has the expected status code and body
	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Contains(t, string(body), payload[1:])
}

// go test -run Test_WithInvalidDataDTO
func Test_WithInvalidDataDTO(t *testing.T) {
	// Create a new app
	app := createAppWithMiddleware2Body()

	// Create a request with wrong body
	payload := `{"name":"computer","price":5000.5,"quantity":"100"}`
	req, _ := http.NewRequest(fiber.MethodPost, endpoint, strings.NewReader(payload))
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(fiber.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

	// Send request to the app
	res, err := app.Test(req)
	require.NoError(t, err)

	// Read the response body into a string
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	// Check that the response has the expected status code and body
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, errorBody, string(body))
}
