package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/youscan/azure-mutex/azuremutex"
	"time"
)

func main() {
	log.Println("Hello")

	mutex := azuremutex.NewMutex("...")
	mutex.Acquire()
	log.Println("Doing some exclusive work")
	time.Sleep(3 * time.Second)
	mutex.Release()
}
