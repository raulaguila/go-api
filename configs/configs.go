package configs

import (
	"crypto/rsa"
	"embed"
	"encoding/base64"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"github.com/raulaguila/go-api/pkg/utils"
)

var (
	//go:embed locales/*
	Locales embed.FS

	AccessPrivateKey  *rsa.PrivateKey
	RefreshPrivateKey *rsa.PrivateKey
)

func init() {
	err := godotenv.Load(path.Join("configs", ".env"))
	if err != nil {
		_, b, _, _ := runtime.Caller(0)
		utils.PanicIfErr(godotenv.Load(path.Join(path.Dir(b), ".env")))
	}

	time.Local, err = time.LoadLocation(os.Getenv("TZ"))
	utils.PanicIfErr(err)

	{
		accessDecodedKey, err := base64.StdEncoding.DecodeString(os.Getenv("ACCESS_TOKEN_PRIVAT"))
		utils.PanicIfErr(err)

		AccessPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(accessDecodedKey)
		utils.PanicIfErr(err)
	}

	{
		refreshDecodedKey, err := base64.StdEncoding.DecodeString(os.Getenv("RFRESH_TOKEN_PRIVAT"))
		utils.PanicIfErr(err)

		RefreshPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(refreshDecodedKey)
		utils.PanicIfErr(err)
	}
}
