package configs

import (
	"embed"
	"os"
	"path"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/raulaguila/go-api/pkg/helper"
)

//go:embed locales/*
var Locales embed.FS

//go:embed version.txt
var version string

func init() {
	err := godotenv.Load(path.Join("configs", ".env"))
	helper.PanicIfErr(err)

	helper.PanicIfErr(os.Setenv("SYS_VERSION", strings.TrimSpace(version)))

	time.Local, err = time.LoadLocation(os.Getenv("TZ"))
	helper.PanicIfErr(err)
}
