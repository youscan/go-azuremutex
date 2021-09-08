package azuremutex

import (
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	acquireInterval      = 10
	leaseDurationSeconds = 60
	renewIntervalSeconds = 10
)

type Locker struct {
	key   string
	mutex *AzureMutex
}

func NewLocker(accountName string, accountKey string, containerName string, key string) *Locker {
	spinLock := Locker{
		key:   key,
		mutex: NewMutex(accountName, accountKey, containerName),
	}
	return &spinLock
}

func (l *Locker) Lock() error {

	err := l.waitLock()
	if err != nil {
		return err
	}

	l.startRenew()

	return nil
}

func (l *Locker) Unlock() error {
	l.stopRenew()
	return l.mutex.Release(l.key)
}

func (l *Locker) waitLock() error {
	for {
		err := l.mutex.Acquire(l.key, leaseDurationSeconds)
		if _, ok := err.(*LeaseAlreadyPresentError); ok {
			log.Debugf("Lock already acquired. Waiting . . .")
			time.Sleep(acquireInterval * time.Second)
			continue
		}
		if err == nil {
			log.Debugf("Locked it!")
		}
		return err
	}
}

func (l *Locker) startRenew() {
	go func() {
		for {
			time.Sleep(renewIntervalSeconds * time.Second)
			l.mutex.Renew(l.key)
			log.Debugf("Lease renewed")
		}
	}()
}

func (l *Locker) stopRenew() {
	// TODO some channels magic for graceful goroutine shutdown
}
