//go:build exclude

package main

import (
	_ "github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/database/pgsql"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/utils"
)

func main() {
	postgresDB := pgsql.ConnectPostgresDB()

	utils.PanicIfErr(postgresDB.AutoMigrate(new(domain.Profile), new(domain.Auth), new(domain.User)))

	utils.PanicIfErr(postgresDB.AutoMigrate(new(domain.Product)))
}
