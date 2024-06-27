package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/i18n"
	"github.com/raulaguila/go-api/pkg/helper"
	"log"
	"mime/multipart"
	"path/filepath"
	"slices"
)

func GetFileFromRequest(formKey string, extensions *[]string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		messages := c.Locals(helper.LocalLang).(*i18n.Translation)
		file, err := c.FormFile(formKey)
		if err != nil || (extensions != nil && !slices.Contains(*extensions, filepath.Ext(file.Filename))) {
			fmt.Println(err.Error())
			return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrInvalidData)
		}

		f, err := file.Open()
		if err != nil {
			return helper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrInvalidData)
		}
		defer func(f multipart.File) {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}(f)

		c.Locals(helper.LocalDTO, &domain.File{
			Name:      file.Filename,
			Extension: filepath.Ext(file.Filename),
			File:      f,
		})
		return c.Next()
	}
}
