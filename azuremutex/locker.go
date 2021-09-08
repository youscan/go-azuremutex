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
	key      string
	mutex    *AzureMutex
	cancel   chan bool
	canceled chan bool
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
	err := l.mutex.Release(l.key)
	if err == nil {
		log.Debugf("Lease released")
	}
	return err
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
	l.cancel = make(chan bool)
	l.canceled = make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(renewIntervalSeconds * time.Second):
				l.mutex.Renew(l.key)
				log.Debugf("Lease renewed")
				break
			case <-l.cancel:
				log.Debugf("Stopping renewing . . .")
				l.canceled <- true
				return
			}
		}
	}()
}

func (l *Locker) stopRenew() {
	l.cancel <- true
	<-l.canceled
	log.Debugf("Stopped")
}
