package handlers

import (
	"github.com/raulaguila/go-api/configs"
	handler2 "github.com/raulaguila/go-api/internal/api/rest/handler"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	service2 "github.com/raulaguila/go-api/internal/api/rest/service"
	"os"
	"strings"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/docs"
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

func initRepositories(db *gorm.DB) {
	profileRepository = repository.NewProfileRepository(db)
	userRepository = repository.NewUserRepository(db)
	productRepository = repository.NewProductRepository(db)
}

func initServices() {
	profileService = service2.NewProfileService(profileRepository)
	userService = service2.NewUserService(userRepository)
	authService = service2.NewAuthService(userRepository)
	productService = service2.NewProductService(productRepository)
}

func initHandlers(app *fiber.App) {
	// Initialize access middlewares
	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, userRepository)
	middleware.MidRefresh = middleware.Auth(configs.RefreshPrivateKey, userRepository)

	// Prepare endpoints for the API.
	handler2.NewMiscHandler(app.Group(""))
	handler2.NewAuthHandler(app.Group("/auth"), authService)
	handler2.NewProfileHandler(app.Group("/profile"), profileService)
	handler2.NewUserHandler(app.Group("/user"), userService)
	handler2.NewProductHandler(app.Group("/product"), productService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		return utils.NewHTTPResponse(c, fiber.StatusNotFound, fiberi18n.MustLocalize(c, "nonExistentRoute"))
	})
}

func HandleRequests(app *fiber.App, db *gorm.DB) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "1" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	initRepositories(db)
	initServices()
	initHandlers(app)

	panic(app.Listen(":" + os.Getenv("API_PORT")))
}
