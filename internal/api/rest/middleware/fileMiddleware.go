package middleware

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"slices"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/pkg/HTTPResponse"
	"github.com/raulaguila/go-api/pkg/utils"
)

type File struct {
	Name      string
	Extension string
	File      io.Reader
}

func GetFileFromRequest(formKey string, extensions *[]string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile(formKey)
		if err != nil || (extensions != nil && !slices.Contains(*extensions, filepath.Ext(file.Filename))) {
			if err != nil {
				fmt.Println(err.Error())
			}
			return HTTPResponse.New(c, fiber.StatusBadRequest, fiberi18n.MustLocalize(c, "invalidData"), nil)
		}

		f, err := file.Open()
		if err != nil {
			return HTTPResponse.New(c, fiber.StatusBadRequest, fiberi18n.MustLocalize(c, "invalidData"), nil)
		}
		defer func(f multipart.File) {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}(f)

		c.Locals(utils.LocalFile, &File{
			Name:      file.Filename,
			Extension: filepath.Ext(file.Filename),
			File:      f,
		})
		return c.Next()
	}
}
