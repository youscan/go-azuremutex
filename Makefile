tests:
	go test -v ./...

lint:
	golangci-lint run
