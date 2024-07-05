package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/raulaguila/go-api/pkg/minioutils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/middleware/language"
	"github.com/raulaguila/go-api/internal/infra/database"
	"github.com/raulaguila/go-api/internal/infra/handlers"
	"github.com/raulaguila/go-api/internal/pkg/i18n"
	"github.com/raulaguila/go-api/pkg/helper"
)

// @title 							Go API
// @description 					This API is a user-friendly solution designed to serve as the foundation for more complex APIs.

// @contact.name					Raul del Aguila
// @contact.email					email@email.com

// @BasePath						/

// @securityDefinitions.apiKey		Bearer
// @in								header
// @name							Authorization
// @description 					Type "Bearer" followed by a space and JWT token.
func main() {
	db, err := database.ConnectPostgresDB()
	helper.PanicIfErr(err)

	minioClient := minioutils.NewMinioClient()
	helper.PanicIfErr(minioClient.InitBucket(context.Background(), os.Getenv("MINIO_BUCKET_FILES"), "*"))

	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		Prefork:               os.Getenv("SYS_PREFORK") == "true",
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		AppName:               "Golang template",
		ReduceMemoryUsage:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return helper.NewHTTPResponse(c, fiber.StatusInternalServerError, err)
		},
		BodyLimit: 50 * 1024 * 1024, // this is the default limit of 50MB
	})

	app.Use(
		recover.New(),
		language.New(language.Config{
			KeyLookup:  "lang",
			ContextKey: helper.LocalLang,
		}),
	)

	if strings.ToLower(os.Getenv("API_LOGGER")) == "true" {
		app.Use(logger.New(logger.Config{
			CustomTags: map[string]logger.LogFunc{
				"xip": func(output logger.Buffer, c *fiber.Ctx, _ *logger.Data, _ string) (int, error) {
					return output.WriteString(fmt.Sprintf("%15s", c.IP()))
				},
				"fullPath": func(output logger.Buffer, c *fiber.Ctx, _ *logger.Data, _ string) (int, error) {
					return output.WriteString(string(c.Request().RequestURI()))
				},
			},
			Format:     "[FIBER:${magenta}${pid}${reset}] ${time} | ${status} | ${latency} | ${xip} | ${method} ${fullPath} ${magenta}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   time.Local.String(),
		}))
	}

	app.Use(
		cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete, fiber.MethodOptions}, ","),
			AllowHeaders:  "*",
			ExposeHeaders: "*",
			MaxAge:        1,
		}),
		limiter.New(limiter.Config{
			Max:        100,
			Expiration: time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				var messages = c.Locals(helper.LocalLang).(*i18n.Translation)
				return helper.NewHTTPResponse(c, fiber.StatusTooManyRequests, messages.ErrManyRequest)
			},
		}),
	)

	handlers.HandleRequests(app, db, minioClient)
}
