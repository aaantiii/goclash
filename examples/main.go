package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/aaantiii/clash"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	creds := make(clash.Credentials)
	emailStr := os.Getenv("EMAILS")
	passwordStr := os.Getenv("PASSWORDS")
	emails := strings.Split(emailStr, ",")
	passwords := strings.Split(passwordStr, ",")

	for i, email := range emails {
		creds[email] = passwords[i]
	}

	_, err := clash.NewClient(creds)
	if err != nil {
		panic(err)
	}
}
