package handlers

import (
	"os"
	"strings"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/docs"
	"github.com/raulaguila/go-api/internal/api/rest/handler"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/api/rest/service"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/repository"
	"github.com/raulaguila/go-api/pkg/utils"
)

var (
	profileRepository domain.ProfileRepository
	userRepository    domain.UserRepository
	productRepository domain.ProductRepository

	profileService domain.ProfileService
	userService    domain.UserService
	authService    domain.AuthService
	productService domain.ProductService
)

func initRepositories(postgresDB *gorm.DB) {
	profileRepository = repository.NewProfileRepository(postgresDB)
	userRepository = repository.NewUserRepository(postgresDB)
	productRepository = repository.NewProductRepository(postgresDB)
}

func initServices() {
	profileService = service.NewProfileService(profileRepository)
	userService = service.NewUserService(userRepository)
	authService = service.NewAuthService(userRepository)
	productService = service.NewProductService(productRepository)
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
	handler.NewProductHandler(app.Group("/product"), productService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		return utils.NewHTTPResponse(c, fiber.StatusNotFound, fiberi18n.MustLocalize(c, "nonExistentRoute"), nil)
	})
}

func HandleRequests(app *fiber.App, postgresDB *gorm.DB) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "1" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	initRepositories(postgresDB)
	initServices()
	initHandlers(app)

	panic(app.Listen(":" + os.Getenv("API_PORT")))
}
