package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aaantiii/goclash"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	keysStr := os.Getenv("KEYS")
	keys := strings.Split(keysStr, ",")

	client, err := goclash.NewWithKeys(keys...)
	if err != nil {
		panic(err)
	}

	// get a player by tag
	player, err := client.GetPlayer("#8QYG8CJ0")
	if err != nil {
		panic(err)
	}
	fmt.Println(player.Name)
}
