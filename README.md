# GoClash
*A blazingly fast, feature rich Clash of Clans API wrapper for Go*

## Installation
To use GoClash, simply run `go get github.com/aaantiii/goclash`.

## Key Features
- **Automatic Key Management** - GoClash automatically manages your API keys, so you don't have to worry about them.
- **Multi Account Support** - GoClash allows you to use multiple API Accounts at once, so that you are not limited to 10 API-Keys.
- **Easy to use** - GoClash is easy to use, and has a very simple API.
- **Caching** - GoClash caches all requests, so that you don't have to worry about rate limits (can be disabled).
- **Concurrency** - GoClash is fully concurrent, so that you can make multiple requests at once.

## Usage
```go
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
```

### More Examples
You can see more examples [here](./examples).
