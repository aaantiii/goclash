package main

import "github.com/aaantiii/goclash"

func main() {
	credentials := goclash.Credentials{"email1": "password1", "email2": "password2"}
	client, err := goclash.New(credentials)
	if err != nil {
		panic(err)
	}

	// get a player by tag
	player, err := client.GetPlayer("#8QYG8CJ0")
	if err != nil {
		panic(err)
	}
	println(player.Name)
}
