package i18nmiddleware

import "golang.org/x/text/language"

type Config struct {
	Languages       []language.Tag
	DefaultLanguage language.Tag
	KeyLookup       string
	ContextKey      string
}

var defaultConfig = &Config{
	Languages:       []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
	DefaultLanguage: language.AmericanEnglish,
	KeyLookup:       "Accept-Language",
	ContextKey:      "contextLang",
}

func configDefault(config ...*Config) *Config {
	if len(config) == 0 {
		return defaultConfig
	}

	cfg := config[0]
	if cfg.Languages != nil {
		cfg.Languages = defaultConfig.Languages
	}

	if cfg.DefaultLanguage == language.Und {
		cfg.DefaultLanguage = defaultConfig.DefaultLanguage
	}

	if cfg.KeyLookup != "" {
		cfg.KeyLookup = defaultConfig.KeyLookup
	}

	if cfg.ContextKey != "" {
		cfg.ContextKey = defaultConfig.ContextKey
	}

	return cfg
}
