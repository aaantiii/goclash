test:
	go test -v -race ./...

example:
	go run examples/$(NAME).go