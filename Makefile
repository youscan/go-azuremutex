tests:
	go test -v ./...

lint:
	golangci-lint run

storage-emulator:
	 docker run -it -p 10000:10000  mcr.microsoft.com/azure-storage/azurite
