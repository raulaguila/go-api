package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"os"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/utils"
)

var (
	MidAccess  fiber.Handler
	MidRefresh fiber.Handler
)

func Auth(base64key string, repo domain.UserRepository) fiber.Handler {
	decodedKey, err := base64.StdEncoding.DecodeString(base64key)
	utils.PanicIfErr(err)

	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(decodedKey)
	utils.PanicIfErr(err)

	return keyauth.New(keyauth.Config{
		KeyLookup:  "header:" + fiber.HeaderAuthorization,
		AuthScheme: "Bearer",
		ContextKey: "token",
		Next: func(c *fiber.Ctx) bool {
			// Filter request to skip middleware
			// true to skip, false to not skip
			c.Locals(utils.LocalUser, new(domain.User))
			return os.Getenv("API_ACCEPT_SKIP_AUTH") == "1" && c.Get("X-Skip-Auth", "false") == "true"
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.NewHTTPResponse(c, fiber.StatusUnauthorized, err.Error())
		},
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			parsedToken, err := jwt.Parse(key, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, err
				}

				return parsedKey, nil
			})
			if err != nil {
				log.Println(err)
				return false, errors.New(fiberi18n.MustLocalize(c, "errGeneric"))
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok || !parsedToken.Valid {
				return false, errors.New("invalid jwt token")
			}

			user, err := repo.GetUserByToken(c.Context(), claims["token"].(string))
			if err != nil {
				log.Println(err)
				return false, errors.New(fiberi18n.MustLocalize(c, "errGeneric"))
			}

			if !user.Auth.Status {
				return false, errors.New(fiberi18n.MustLocalize(c, "disabledUser"))
			}

			c.Locals(utils.LocalUser, user)
			return true, nil
		},
	})
}
