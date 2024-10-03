package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

func pgConnect(dbName string) *gorm.DB {
	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%v port=%s sslmode=disable TimeZone=%v", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), dbName, os.Getenv("POSTGRES_PORT"), time.Local.String())
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		NowFunc:     time.Now,
		PrepareStmt: true,
	})
	helper.PanicIfErr(err)

	return db
}

func createDataBase() {
	db := pgConnect("postgres")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, err := db.WithContext(ctx).DB()
	helper.PanicIfErr(err)
	defer func(con *sql.DB) {
		_ = con.Close()
	}(con)

	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %v;", os.Getenv("POSTGRES_BASE"))).Error; err != nil {
		switch {
		case errors.Is(pgutils.HandlerError(err), pgutils.ErrDatabaseAlreadyExists):
		default:
			helper.PanicIfErr(err)
		}
	}
}

func ConnectPostgresDB() (*gorm.DB, error) {
	createDataBase()

	db := pgConnect(os.Getenv("POSTGRES_BASE"))
	db = db.WithContext(context.Background())
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS unaccent;").Error; err != nil {
		return nil, err
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, err
	}

	autoMigrate(db)
	go createDefaults(db)
	return db, nil
}
