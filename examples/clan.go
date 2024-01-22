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

	client, err := clash.New(creds)
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
	clans, err := client.SearchClans(clash.SearchClanParams{Name: "Lost 5", PagingParams: &clash.PagingParams{Limit: 5}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", clans)
}
