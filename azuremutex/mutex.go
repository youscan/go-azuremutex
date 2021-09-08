package azuremutex

import (
	log "github.com/sirupsen/logrus"
)

type AzureMutex interface {
	Acquire()
	Release()
}

type mutexState struct {
	connectionString string
}

func NewMutex(connectionString string) AzureMutex {
	return mutexState{
		connectionString: connectionString,
	}
}

func (m mutexState) Acquire() {
	log.Println("Acquire not implemented")
}

func (m mutexState) Release() {
	log.Println("Release not implemented")
}
