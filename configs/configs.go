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

	"github.com/raulaguila/go-api/pkg/packhub"
	"github.com/raulaguila/go-api/pkg/utils"
)

var (
	//go:embed locales/*
	Locales embed.FS

	AccessPrivateKey *rsa.PrivateKey
	AccessExpiration time.Duration

	RefreshPrivateKey *rsa.PrivateKey
	RefreshExpiration time.Duration
)

func init() {
	err := godotenv.Load(path.Join("configs", ".env"))
	if err != nil {
		_, b, _, _ := runtime.Caller(0)
		packhub.PanicIfErr(godotenv.Load(path.Join(path.Dir(b), ".env")))
	}

	time.Local, err = time.LoadLocation(os.Getenv("TZ"))
	packhub.PanicIfErr(err)

	{
		accessDecodedKey, err := base64.StdEncoding.DecodeString(os.Getenv("ACCESS_TOKEN"))
		packhub.PanicIfErr(err)

		AccessPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(accessDecodedKey)
		packhub.PanicIfErr(err)

		AccessExpiration, err = utils.DurationFromString(os.Getenv("ACCESS_TOKEN_EXPIRE"), time.Minute)
		packhub.PanicIfErr(err)
	}

	{
		refreshDecodedKey, err := base64.StdEncoding.DecodeString(os.Getenv("RFRESH_TOKEN"))
		packhub.PanicIfErr(err)

		RefreshPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(refreshDecodedKey)
		packhub.PanicIfErr(err)

		RefreshExpiration, err = utils.DurationFromString(os.Getenv("RFRESH_TOKEN_EXPIRE"), time.Minute)
		packhub.PanicIfErr(err)
	}
}
