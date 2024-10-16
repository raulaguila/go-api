package i18nmiddleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

func New(config ...*Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c *fiber.Ctx) error {
		lang, err := language.Parse(c.Get(cfg.KeyLookup, cfg.DefaultLanguage.String()))
		if err != nil || !slices.Contains(cfg.Languages, lang) {
			lang = cfg.DefaultLanguage
		}

		c.Locals(cfg.ContextKey, lang)
		return c.Next()
	}
}
