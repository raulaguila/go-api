package database

import (
	"context"
	"github.com/raulaguila/go-api/pkg/utils"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/helper"
)

// autoMigrate migrates the database schema for Profile, Auth, User, and Product models using GORM's AutoMigrate function.
func autoMigrate(db *gorm.DB) {
	helper.PanicIfErr(db.AutoMigrate(new(domain.Profile)))
	helper.PanicIfErr(db.AutoMigrate(new(domain.Auth)))
	helper.PanicIfErr(db.AutoMigrate(new(domain.User)))

	helper.PanicIfErr(db.AutoMigrate(new(domain.Product)))
}

// createDefaults initializes the database with a default profile and user if they do not already exist.
func createDefaults(db *gorm.DB) {
	profile := &domain.Profile{
		Name: "ROOT",
		Permissions: map[string]any{
			"user":    true,
			"product": true,
			"profile": true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	helper.PanicIfErr(db.WithContext(ctx).FirstOrCreate(profile, "name = ?", profile.Name).Error)

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADM_PASS")), bcrypt.DefaultCost)
	helper.PanicIfErr(err)

	user := &domain.User{
		Name:  os.Getenv("ADM_NAME"),
		Email: os.Getenv("ADM_MAIL"),
		Auth: &domain.Auth{
			Status:    true,
			ProfileID: profile.ID,
			Token:     utils.Pointer(uuid.New().String()),
			Password:  utils.Pointer(string(hash)),
		},
	}

	helper.PanicIfErr(db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).FirstOrCreate(user, "mail = ?", user.Email).Error)
}
