package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/aaantiii/goclash"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	creds := make(goclash.Credentials)
	emailStr := os.Getenv("EMAILS")
	passwordStr := os.Getenv("PASSWORDS")
	emails := strings.Split(emailStr, ",")
	passwords := strings.Split(passwordStr, ",")

	for i, email := range emails {
		creds[email] = passwords[i]
	}

	client, err := goclash.New(creds)
	if err != nil {
		panic(err)
	}

	// get a clan by tag
	clan, err := client.GetClan("#2QC0QQPQ2")
	if err != nil {
		panic(err)
	}
	fmt.Println(clan.Name)

	// get clan war log
	log, err := client.GetClanWarLog("#2QC0QQPQ2", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", log)

	// search clans
	clans, err := client.SearchClans(goclash.SearchClanParams{Name: "LOST", PagingParams: &goclash.PagingParams{Limit: 5}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", clans)
}
