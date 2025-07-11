//go:build exclude

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/nexidian/gocliselect"
	"golang.org/x/crypto/bcrypt"

	"github.com/raulaguila/go-api/pkg/packhub"
)

func generateRSAPrivateToken() {
	fmt.Println()
	menu := gocliselect.NewMenu("Bit size of the key")

	menu.AddItem("1024 bits", "1024")
	menu.AddItem("2048 bits", "2048")
	menu.AddItem("3072 bits", "3072")
	menu.AddItem("4096 bits", "4096")

	bits, err := strconv.Atoi(menu.Display())
	packhub.PanicIfErr(err)

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	packhub.PanicIfErr(err)

	fmt.Printf("\nPrivate key: %v\n\n", base64.StdEncoding.EncodeToString(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})))
}

func hashPassword() {
	fmt.Print("\nUser password: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	hash, err := bcrypt.GenerateFromPassword([]byte(scanner.Text()), bcrypt.DefaultCost)
	packhub.PanicIfErr(err)

	fmt.Printf("Hash: %s\n\n", hash)
}

func generateUserToken() {
	fmt.Printf("\nUser token: %s\n\n", uuid.New().String())
}

func main() {
	printMenu := func() string {
		menu := gocliselect.NewMenu("Chose an option")

		menu.AddItem("Generate RSA Token", "rsa")
		menu.AddItem("Hash user password", "hash")
		menu.AddItem("New user token", "token")
		menu.AddItem("Exit", "exit")

		return menu.Display()
	}

	for {
		switch choice := printMenu(); choice {
		case "exit":
			return
		case "rsa":
			generateRSAPrivateToken()
		case "token":
			generateUserToken()
		case "hash":
			hashPassword()
		}
	}
}
