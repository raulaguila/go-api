package handler

import (
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"

	"github.com/raulaguila/go-api/configs"
)

func setupMiscApp() *fiber.App {
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
	NewMiscHandler(app.Group("/"))

	return app
}

func TestMiscHandler_healthCheck(t *testing.T) {
	tests := []generalHandlerTest{
		{
			name:         "success",
			method:       fiber.MethodGet,
			endpoint:     "/",
			body:         nil,
			setupMocks:   func() {},
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "error",
			method:       fiber.MethodGet,
			endpoint:     "/invalid",
			body:         nil,
			setupMocks:   func() {},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupMiscApp())
}
