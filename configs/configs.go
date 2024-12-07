package configs

import (
	"embed"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/raulaguila/go-api/pkg/utils"
)

//go:embed locales/*
var Locales embed.FS

//go:embed version.txt
var version string

// init initializes the environment configuration by loading variables from a .env file and setting system properties.
// It panics if any error occurs during the loading of environment variables or setting time location.
func init() {
	err := godotenv.Load(path.Join("configs", ".env"))
	if err != nil {
		_, b, _, _ := runtime.Caller(0)
		utils.PanicIfErr(godotenv.Load(path.Join(path.Dir(b), ".env")))
	}

	utils.PanicIfErr(os.Setenv("SYS_VERSION", strings.TrimSpace(version)))

	time.Local, err = time.LoadLocation(os.Getenv("TZ"))
	utils.PanicIfErr(err)
}
