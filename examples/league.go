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

	// get all leagues (home village)
	leagues, err := client.GetLeagues(nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", leagues)

	// get legend league seasons, limit to 5 results
	seasons, err := client.GetLeagueSeasons(goclash.LeagueLegend, &goclash.PagingParams{Limit: 5})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", seasons)
}
