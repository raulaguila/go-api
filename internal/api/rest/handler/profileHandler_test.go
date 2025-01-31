package handler

import (
	"fmt"
	"golang.org/x/text/language"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/internal/pkg/repository"
	"github.com/raulaguila/go-api/internal/pkg/service"
	"github.com/raulaguila/go-api/pkg/utils"
)

func setupProfileApp(db *gorm.DB) *fiber.App {
	repo := repository.NewProfileRepository(db)

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
	NewProfileHandler(app.Group(""), service.NewProfileService(repo))

	return app
}

func TestProfileHandler_getProfiles(t *testing.T) {
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

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))

				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 02", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 03", Permissions: pq.StringArray{"read"}}).Error)

				return setupProfileApp(db)
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

func TestProfileHandler_getProfile(t *testing.T) {
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

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))

				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 02", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 03", Permissions: pq.StringArray{"read"}}).Error)

				return setupProfileApp(db)
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

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))

				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 02", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 03", Permissions: pq.StringArray{"read"}}).Error)

				return setupProfileApp(db)
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

func TestProfileHandler_createProfile(t *testing.T) {
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Profile 01", "permissions": ["read"]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				return setupProfileApp(db)
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "duplicate Profile",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Profile 01", "permissions": ["read"]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				return setupProfileApp(db)
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

func TestProfileHandler_updateProfile(t *testing.T) {
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func() *fiber.App
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Profile 01", "permissions": ["read"]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				return setupProfileApp(db)
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "duplicate Profile",
			endpoint: "/2",
			body:     strings.NewReader(`{"name": "Profile 01", "permissions": ["read"]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 02", Permissions: pq.StringArray{"read"}}).Error)
				return setupProfileApp(db)
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Profile 01", "permissions": ["read"]}`),
			setup: func() *fiber.App {
				db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				utils.PanicIfErr(err)

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				return setupProfileApp(db)
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

func TestProfileHandler_deleteProfile(t *testing.T) {
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

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 02", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 03", Permissions: pq.StringArray{"read"}}).Error)
				return setupProfileApp(db)
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

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				return setupProfileApp(db)
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

				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 01", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 02", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 03", Permissions: pq.StringArray{"read"}}).Error)
				return setupProfileApp(db)
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
