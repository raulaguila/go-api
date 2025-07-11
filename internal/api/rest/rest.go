package rest

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
	"github.com/gofiber/swagger"
	"github.com/minio/minio-go/v7"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/docs"
	"github.com/raulaguila/go-api/internal/api/rest/handler"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/HTTPResponse"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/repository"
	"github.com/raulaguila/go-api/internal/pkg/service"
	"github.com/raulaguila/go-api/pkg/packhub"
)

var (
	profileRepository domain.ProfileRepository
	userRepository    domain.UserRepository

	authService    domain.AuthService
	profileService domain.ProfileService
	userService    domain.UserService
)

func initRepositories(postgresDB *gorm.DB, minioClient *minio.Client) {
	profileRepository = repository.NewProfileRepository(postgresDB)
	userRepository = repository.NewUserRepository(postgresDB)
}

func initServices() {
	profileService = service.NewProfileService(profileRepository)
	authService = service.NewAuthService(userRepository)
	userService = service.NewUserService(userRepository)
}

func initHandlers(app *fiber.App) {
	// Initialize access middlewares
	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, userRepository)
	middleware.MidRefresh = middleware.Auth(configs.RefreshPrivateKey, userRepository)

	// Prepare endpoints for the API.
	handler.NewMiscHandler(app.Group(""))
	handler.NewAuthHandler(app.Group("/auth"), authService)

	handler.NewProfileHandler(app.Group("/profile"), profileService)

	handler.NewUserHandler(app.Group("/user"), userService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		return HTTPResponse.New(c, fiber.StatusNotFound, fiberi18n.MustLocalize(c, "nonExistentRoute"), nil)
	})
}

func start(app *fiber.App, postgresDB *gorm.DB, minioClient *minio.Client) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "1" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	initRepositories(postgresDB, minioClient)
	initServices()
	initHandlers(app)

	packhub.PanicIfErr(app.Listen(":" + os.Getenv("API_PORT")))
}

func New(postgresDB *gorm.DB, minioClient *minio.Client) {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		Prefork:               os.Getenv("API_ENABLE_PREFORK") == "1",
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		AppName:               "API Backend",
		ReduceMemoryUsage:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return HTTPResponse.New(c, fiber.StatusInternalServerError, err.Error(), nil)
		},
		BodyLimit: 100 * 1024 * 1024,
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
				"xauth": func(output logger.Buffer, c *fiber.Ctx, _ *logger.Data, _ string) (int, error) {
					if key := c.Get("Authorization", ""); key != "" {
						return output.WriteString(fmt.Sprintf(":%s", strings.TrimPrefix(key, "Bearer ")[:min(len(strings.TrimPrefix(key, "Bearer ")), 10)]))
					}
					return output.WriteString("")
				},
			},
			Format:     "[FIBER:${magenta}${xid}${reset}] ${time} | ${status} | ${latency} | ${xip} | ${method} ${fullPath} ${yellow}\"${reqHeader:Accept-Language}${xauth}\"${reset} ${magenta}${error}${reset}\n",
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
			MaxAge:        -1,
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
			Next:       nil,
			Max:        300,
			Expiration: 30 * time.Second,
			LimitReached: func(c *fiber.Ctx) error {
				return HTTPResponse.New(c, fiber.StatusTooManyRequests, fiberi18n.MustLocalize(c, "manyRequests"), nil)
			},
		}),
	)

	start(app, postgresDB, minioClient)
}
