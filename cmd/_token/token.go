//go:build exclude

package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/raulaguila/go-api/pkg/utils"
)

// main is the entry point of the application. It generates a bcrypt hash for a given password and prints the password and hash.
func main() {
	password := "12345678"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	utils.PanicIfErr(err)

	fmt.Printf("Password: %s\nHash: %s\n", password, hash)
}
