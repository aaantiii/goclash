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

	// get a player by tag
	player, err := client.GetPlayer("#8QYG8CJ0")
	if err != nil {
		panic(err)
	}
	fmt.Println(player.Name)

	// concurrently get multiple players by tag
	players, err := client.GetPlayersWithError("#8QYG8CJ0", "#Q9RY8YRYJ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", players)

	// verify player with api token
	verification, err := client.VerifyPlayer("#8QYG8CJ0", "apiToken")
	if err != nil {
		panic(err)
	}
	if verification.IsOk() {
		fmt.Println("verification successful")
	}
}
