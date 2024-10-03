package handlers

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/docs"
	"github.com/raulaguila/go-api/internal/api/handler"
	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/api/service"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/i18n"
	"github.com/raulaguila/go-api/internal/pkg/repository"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/minioutils"
)

var (
	profileRepository    domain.ProfileRepository
	userRepository       domain.UserRepository
	departmentRepository domain.DepartmentRepository

	profileService    domain.ProfileService
	userService       domain.UserService
	authService       domain.AuthService
	departmentService domain.DepartmentService
)

// initRepositories Initialize all repositories.
func initRepositories(db *gorm.DB, minioClient *minioutils.Minio) {
	profileRepository = repository.NewProfileRepository(db)
	userRepository = repository.NewUserRepository(db, minioClient)
	departmentRepository = repository.NewDepartmentRepository(db)
}

// initServices Initialize all services.
func initServices() {
	profileService = service.NewProfileService(profileRepository)
	userService = service.NewUserService(userRepository)
	authService = service.NewAuthService(userRepository)
	departmentService = service.NewDepartmentService(departmentRepository)
}

func initHandlers(app *fiber.App) {
	// Initialize access middlewares
	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), userRepository)
	middleware.MidRefresh = middleware.Auth(os.Getenv("RFRESH_TOKEN_PUBLIC"), userRepository)

	// Prepare endpoints for the API.
	handler.NewMiscHandler(app.Group(""))
	handler.NewAuthHandler(app.Group("/auth"), authService)
	handler.NewProfileHandler(app.Group("/profile"), profileService)
	handler.NewUserHandler(app.Group("/user"), userService)
	handler.NewDepartmentHandler(app.Group("/department"), departmentService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		messages := c.Locals(helper.LocalLang).(*i18n.Translation)
		return helper.NewHTTPResponse(c, fiber.StatusNotFound, messages.ErrorNonexistentRoute)
	})
}

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
