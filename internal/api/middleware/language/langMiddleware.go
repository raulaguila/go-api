package language

import (
	"os"
	"slices"
	"strings"

	myi18n "github.com/raulaguila/go-api/internal/pkg/i18n"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	KeyLookup  string
	ContextKey string
}

var defaultConfig = Config{
	KeyLookup:  "lang",
	ContextKey: "localLang",
}

func New(config ...Config) fiber.Handler {
	cfg := defaultConfig
	if len(config) > 0 {
		aux := config[0]
		if aux.KeyLookup != "" {
			cfg.KeyLookup = aux.KeyLookup
		}
		if aux.ContextKey != "" {
			cfg.ContextKey = aux.ContextKey
		}
	}

	return func(c *fiber.Ctx) error {
		lang := strings.ToLower(c.Query(cfg.KeyLookup, os.Getenv("SYS_LANGUAGE")))[:2]

		if !slices.Contains(strings.Split(os.Getenv("SYS_LANGUAGES"), ","), lang) {
			lang = os.Getenv("SYS_LANGUAGE")
		}

		c.Locals(cfg.ContextKey, myi18n.TranslationsI18n[lang])
		return c.Next()
	}
}
