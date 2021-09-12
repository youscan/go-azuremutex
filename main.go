package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/youscan/azure-mutex/azuremutex"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {

	options := azuremutex.MutexOptions{
		ContainerName:      "locks",
		UseStorageEmulator: true,
	}
	lock := azuremutex.NewLocker(options, "test")
	log.Println("Waiting for a lock . . .")
	err := lock.Lock()
	panicWhenError(err)

	for i := 0; i < 2; i++ {
		log.Infof("Doing some exclusive work #%d", i)
		time.Sleep(15 * time.Second)
	}
	log.Info("Work done")

	log.Println("Releasing lock")
	err = lock.Unlock()
	panicWhenError(err)
}

func panicWhenError(err error) {
	if err != nil {
		log.Fatalf("Can't continue: %v", err)
		log.Exit(-1)
	}
}
