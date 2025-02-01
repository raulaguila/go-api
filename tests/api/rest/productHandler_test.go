package rest

import (
	"fmt"
	"golang.org/x/text/language"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest/handler"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/internal/pkg/repository"
	"github.com/raulaguila/go-api/internal/pkg/service"
	"github.com/raulaguila/go-api/pkg/utils"
)

func setupProductApp(db *gorm.DB) *fiber.App {
	repo := repository.NewProductRepository(db)

	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, &mocks.UserRepositoryMock{})

	app := fiber.New()
	app.Use(fiberi18n.New(&fiberi18n.Config{
		Next: func(c *fiber.Ctx) bool {
			return false
		},
		RootPath:        "./locales",
		AcceptLanguages: []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
		DefaultLanguage: language.AmericanEnglish,
		Loader:          &fiberi18n.EmbedLoader{FS: configs.Locales},
	}))
	handler.NewProductHandler(app.Group(""), service.NewProductService(repo))

	return app
}

func TestProductHandler_getProducts(t *testing.T) {
	tests := []struct {
		name, endpoint string
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))

				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 02"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 03"}).Error)

				return setupProductApp(db)
			},
			expectedCode: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, tt.endpoint, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := tt.setup().Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_getProduct(t *testing.T) {
	tests := []struct {
		name, endpoint string
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))

				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 02"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 03"}).Error)

				return setupProductApp(db)
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "not found",
			endpoint: "/500",
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))

				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 02"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 03"}).Error)

				return setupProductApp(db)
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, tt.endpoint, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := tt.setup().Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_createProduct(t *testing.T) {
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "duplicate product",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, tt.endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := tt.setup().Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_updateProduct(t *testing.T) {
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "duplicate product",
			endpoint: "/2",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 02"}).Error)
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			req := httptest.NewRequest(fiber.MethodPut, tt.endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := tt.setup().Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_deleteProduct(t *testing.T) {
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 02"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 03"}).Error)
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "not found",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:     "invalid id",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": ["a", "b", "c", "d"]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 01"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 02"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 03"}).Error)
				return setupProductApp(db)
			},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			req := httptest.NewRequest(fiber.MethodDelete, tt.endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := tt.setup().Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}
