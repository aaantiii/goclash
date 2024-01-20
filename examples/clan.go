package main

import (
	"fmt"
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

	client, err := clash.NewClient(creds)
	if err != nil {
		panic(err)
	}

	// get a clan by tag
	clan, err := client.GetClan("#2QC0QQPQ2")
	if err != nil {
		panic(err)
	}
	fmt.Println(clan.Name)

	// concurrently get multiple clans by tag
	clans, err := client.GetClans("#2QC0QQPQ2", "#2820UPPQC", "#2LG222Q0L", "#2YVJV8VC0")
	if err != nil {
		panic(err)
	}
	fmt.Println(clans.Tags())
}
