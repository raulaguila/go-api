//go:build exclude

package main

import (
	_ "github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/infra/database"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/utils"
)

func main() {
	db := database.ConnectPostgresDB()

	utils.PanicIfErr(db.AutoMigrate(new(domain.Profile)))
	utils.PanicIfErr(db.AutoMigrate(new(domain.Auth)))
	utils.PanicIfErr(db.AutoMigrate(new(domain.User)))

	utils.PanicIfErr(db.AutoMigrate(new(domain.Product)))
}
