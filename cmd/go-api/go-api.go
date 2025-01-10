package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/text/language"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/infra/database"
	"github.com/raulaguila/go-api/internal/infra/handlers"
	"github.com/raulaguila/go-api/pkg/utils"
)

// @title 							Go API
// @description 					This API is a user-friendly solution designed to serve as the foundation for more complex APIs.

// @contact.name					Raul del Aguila
// @contact.email					email@email.com

// @BasePath						/

// @securityDefinitions.apiKey		Bearer
// @in								header
// @name							Authorization
// @description 					Type "Bearer" followed by a space and the JWT token.
func main() {
	db := database.ConnectPostgresDB()

	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		Prefork:               os.Getenv("API_ENABLE_PREFORK") == "1",
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		AppName:               "Golang template",
		ReduceMemoryUsage:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.NewHTTPResponse(c, fiber.StatusInternalServerError, err.Error())
		},
		BodyLimit: 50 * 1024 * 1024,
	})

	app.Use(recover.New())

	if strings.ToLower(os.Getenv("API_LOGGER")) == "1" {
		app.Use(logger.New(logger.Config{
			CustomTags: map[string]logger.LogFunc{
				"xid": func(output logger.Buffer, _ *fiber.Ctx, data *logger.Data, _ string) (int, error) {
					return output.WriteString(fmt.Sprintf("%6s", data.Pid))
				},
				"fullPath": func(output logger.Buffer, c *fiber.Ctx, _ *logger.Data, _ string) (int, error) {
					return output.WriteString(c.OriginalURL())
				},
				"xip": func(output logger.Buffer, c *fiber.Ctx, _ *logger.Data, _ string) (int, error) {
					return output.WriteString(fmt.Sprintf("%15s", c.IP()))
				},
			},
			Format:     "[FIBER:${magenta}${xid}${reset}] ${time} | ${status} | ${latency} | ${xip} | ${method} ${fullPath} ${yellow}\"${reqHeader:Accept-Language}\"${reset} ${magenta}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   time.Local.String(),
		}))
	}

	app.Use(
		cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions}, ","),
			AllowHeaders:  "*",
			ExposeHeaders: "*",
			MaxAge:        1,
		}),
		fiberi18n.New(&fiberi18n.Config{
			Next: func(c *fiber.Ctx) bool {
				return false
			},
			RootPath:        "./locales",
			AcceptLanguages: []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
			DefaultLanguage: language.AmericanEnglish,
			Loader:          &fiberi18n.EmbedLoader{FS: configs.Locales},
		}),
		limiter.New(limiter.Config{
			Max:        100,
			Expiration: time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return utils.NewHTTPResponse(c, fiber.StatusTooManyRequests, fiberi18n.MustLocalize(c, "manyRequests"))
			},
		}),
	)

	handlers.HandleRequests(app, db)
}
