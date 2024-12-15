package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/pkg/utils"
)

func ConnectPostgresDB() *gorm.DB {
	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%v port=%s sslmode=disable TimeZone=%v", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_BASE"), os.Getenv("POSTGRES_PORT"), time.Local.String())
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		NowFunc:     time.Now,
		PrepareStmt: true,
	})
	utils.PanicIfErr(err)

	return db
}
