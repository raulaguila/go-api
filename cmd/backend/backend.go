package main

import (
	_ "github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest"
	"github.com/raulaguila/go-api/internal/infra/minio"
	"github.com/raulaguila/go-api/internal/infra/pgsql"
)

// @title 							Go API
// @version							1.0.0
// @description 					This API is a user-friendly solution designed to serve as the foundation for more complex APIs.

// @contact.name					Raul del Aguila
// @contact.url 					https://github.com/raulaguila
// @contact.email					email@email.com

// @BasePath						/

// @securityDefinitions.apiKey		Bearer
// @in								header
// @name							Authorization
// @description 					Type "Bearer" followed by a space and the JWT token.
func main() {
	rest.New(pgsql.ConnectPostgresDB(), minio.ConnectMinio())
}
