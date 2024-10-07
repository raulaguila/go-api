package language

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Languages       []string
	DefaultLanguage string
	KeyLookup       string
	ContextKey      string
}

var defaultConfig = Config{
	Languages:       []string{"en-US", "pt-BR"},
	DefaultLanguage: "en-US",
	KeyLookup:       "Accept-Language",
	ContextKey:      "localLang",
}

func New(config ...Config) fiber.Handler {
	cfg := defaultConfig
	if len(config) > 0 {
		aux := config[0]
		if aux.Languages != nil {
			cfg.Languages = aux.Languages
		}
		for {
			if index := slices.Index(cfg.Languages, ""); index != -1 {
				cfg.Languages = append(cfg.Languages[:index], cfg.Languages[index+1:]...)
				continue
			}
			break
		}
		if aux.DefaultLanguage != "" {
			cfg.DefaultLanguage = aux.DefaultLanguage
		}
		if aux.KeyLookup != "" {
			cfg.KeyLookup = aux.KeyLookup
		}
		if aux.ContextKey != "" {
			cfg.ContextKey = aux.ContextKey
		}
	}

	return func(c *fiber.Ctx) error {
		lang := strings.ToLower(c.Get(cfg.KeyLookup, cfg.DefaultLanguage))
		c.Set(cfg.KeyLookup, fmt.Sprintf("%v;", strings.Join(cfg.Languages, ";")))

		if !slices.Contains(cfg.Languages, lang) {
			lang = cfg.DefaultLanguage
		}

		c.Locals(cfg.ContextKey, lang)
		return c.Next()
	}
}
