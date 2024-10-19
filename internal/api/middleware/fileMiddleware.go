package middleware

import (
	"fmt"
	"github.com/gofiber/contrib/fiberi18n/v2"

	"log"
	"mime/multipart"
	"path/filepath"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/helper"
)

func GetFileFromRequest(formKey string, extensions *[]string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile(formKey)
		if err != nil || (extensions != nil && !slices.Contains(*extensions, filepath.Ext(file.Filename))) {
			if err != nil {
				fmt.Println(err.Error())
			}
			return helper.NewHTTPResponse(c, fiber.StatusBadRequest, fiberi18n.MustLocalize(c, "invalidData"))
		}

		f, err := file.Open()
		if err != nil {
			return helper.NewHTTPResponse(c, fiber.StatusBadRequest, fiberi18n.MustLocalize(c, "invalidData"))
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
