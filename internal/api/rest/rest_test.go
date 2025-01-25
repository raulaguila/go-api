package rest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMiddleware(t *testing.T) {
	app := fiber.New()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	os.Setenv("API_PORT", "3000")
	defer os.Unsetenv("API_PORT")

	go func() {
		go func() {
			New(db)
		}()
		time.Sleep(12 * time.Second) // Let the server initialize
	}()

	tests := []struct {
		name       string
		path       string
		method     string
		statusCode int
	}{
		{
			name:       "not found route",
			path:       "/undefined-route",
			method:     http.MethodGet,
			statusCode: fiber.StatusNotFound,
		},
	}

	time.Sleep(10 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp, err := app.Test(req, -500)

			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}
