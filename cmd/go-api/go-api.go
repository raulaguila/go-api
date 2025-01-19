package main

import (
	_ "github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest"
	"github.com/raulaguila/go-api/internal/database/pgsql"
)

// @title 							Go API
// @description 					This API is a user-friendly solution designed to serve as the foundation for more complex APIs.

// @contact.name					Raul del Aguila
// @contact.email					email@email.com

// @BasePath						/

// @securityDefinitions.apiKey		Bearer
// @in								header
// @name							Authorization
// @description 					Type "Bearer" followed by a space and the JWT token.
func main() {
	postgresDB := pgsql.ConnectPostgresDB()

	rest.New(postgresDB)
}
