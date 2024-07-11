package database

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/helper"
)

func autoMigrate(db *gorm.DB) {
	helper.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
	helper.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
	helper.PanicIfErr(db.AutoMigrate(&domain.User{}))

	helper.PanicIfErr(db.AutoMigrate(&domain.Department{}))
}

func createDefaults(db *gorm.DB) {
	profile := &domain.Profile{
		Name: "ROOT",
		Permissions: map[string]any{
			"user":       true,
			"department": true,
			"profile":    true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	helper.PanicIfErr(db.WithContext(ctx).FirstOrCreate(profile, "name = ?", profile.Name).Error)

	token := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADM_PASS")), bcrypt.DefaultCost)
	helper.PanicIfErr(err)
	pass := string(hash)

	user := &domain.User{
		Name:  os.Getenv("ADM_NAME"),
		Email: os.Getenv("ADM_MAIL"),
		Auth: &domain.Auth{
			Status:    true,
			ProfileID: profile.ID,
			Token:     &token,
			Password:  &pass,
		},
	}

	helper.PanicIfErr(db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).FirstOrCreate(user, "mail = ?", user.Email).Error)
}
