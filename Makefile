tests: start-storage-emulator
	go test -race -v ./...

lint:
	golangci-lint run

start-storage-emulator:
	docker run -d --rm --name azurite -p 10000:10000  mcr.microsoft.com/azure-storage/azurite

stop-storage-emulator:
	docker stop azurite

# brew install act
test-github-actions:
	act
