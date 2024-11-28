package handlers

import (
	"os"
	"strings"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/docs"
	"github.com/raulaguila/go-api/internal/api/handler"
	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/api/service"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/repository"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/minioutils"
)

// profileRepository is an instance of the ProfileRepository interface for managing profile data operations.
// userRepository is an instance of the UserRepository interface for managing user data operations.
// productRepository is an instance of the ProductRepository interface for managing product data operations.
// profileService is an instance of the ProfileService interface for handling profile-related business logic.
// userService is an instance of the UserService interface for handling user-related business logic.
// authService is an instance of the AuthService interface for managing authentication processes.
// productService is an instance of the ProductService interface for handling product-related business logic.
var (
	profileRepository domain.ProfileRepository
	userRepository    domain.UserRepository
	productRepository domain.ProductRepository

	profileService domain.ProfileService
	userService    domain.UserService
	authService    domain.AuthService
	productService domain.ProductService
)

// initRepositories initializes the repositories for profile, user, and product entities using the provided database
// connection and MinIO client. It assigns the created repository instances to their respective global variables.
func initRepositories(db *gorm.DB, minioClient *minioutils.Minio) {
	profileRepository = repository.NewProfileRepository(db)
	userRepository = repository.NewUserRepository(db, minioClient)
	productRepository = repository.NewProductRepository(db)
}

// initServices initializes all necessary services for the application by wiring up repositories to service instances.
func initServices() {
	profileService = service.NewProfileService(profileRepository)
	userService = service.NewUserService(userRepository)
	authService = service.NewAuthService(userRepository)
	productService = service.NewProductService(productRepository)
}

// initHandlers initializes all route handlers and middleware for the given Fiber application instance.
func initHandlers(app *fiber.App) {
	// Initialize access middlewares
	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), userRepository)
	middleware.MidRefresh = middleware.Auth(os.Getenv("RFRESH_TOKEN_PUBLIC"), userRepository)

	// Prepare endpoints for the API.
	handler.NewMiscHandler(app.Group(""))
	handler.NewAuthHandler(app.Group("/auth"), authService)
	handler.NewProfileHandler(app.Group("/profile"), profileService)
	handler.NewUserHandler(app.Group("/user"), userService)
	handler.NewProductHandler(app.Group("/product"), productService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		return helper.NewHTTPResponse(c, fiber.StatusNotFound, fiberi18n.MustLocalize(c, "nonExistentRoute"))
	})
}

// HandleRequests configures route handlers for the app, initializes dependencies and starts the server.
func HandleRequests(app *fiber.App, db *gorm.DB, minioClient *minioutils.Minio) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "1" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	initRepositories(db, minioClient)
	initServices()
	initHandlers(app)

	panic(app.Listen(":" + os.Getenv("API_PORT")))
}
