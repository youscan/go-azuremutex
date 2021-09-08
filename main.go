package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/youscan/azure-mutex/azuremutex"
	"time"
)

func main() {
	const (
		accountName = "*****"
		accountKey  = "*****"
		container   = "locks"
	)

	log.Println("Acquiring mutex")

	mutex := azuremutex.NewMutex(accountName, accountKey, container)
	err := mutex.Acquire("test")
	panicWhenError(err)

	log.Println("Doing some exclusive work")
	time.Sleep(1 * time.Second)

	log.Println("Renewing lock")
	err = mutex.Renew("test")
	panicWhenError(err)
	time.Sleep(10 * time.Second)

	log.Println("Releasing lock")
	err = mutex.Release("test")
	panicWhenError(err)
}

func panicWhenError(err error) {
	if err != nil {
		log.Fatalf("Can't continue: %v", err)
		log.Exit(-1)
	}
}
