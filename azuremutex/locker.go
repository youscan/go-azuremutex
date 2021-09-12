package azuremutex

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	acquireInterval      = 10
	leaseDurationSeconds = 60
	renewIntervalSeconds = 10
)

type Locker struct {
	key            string
	mutex          *AzureMutex
	cancelContext  context.CancelFunc
	cancelRequired chan *sync.WaitGroup
}

func NewLocker(options MutexOptions, key string) *Locker {
	ctx, cancel := context.WithCancel(context.Background())
	spinLock := Locker{
		key:           key,
		cancelContext: cancel,
		mutex:         NewMutexWithContext(options, ctx),
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
	defer l.cancelContext()
	log.Debugf("Unlocked")
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
	l.cancelRequired = make(chan *sync.WaitGroup)
	go func() {
		var wg *sync.WaitGroup
		for {
			select {
			case <-time.After(renewIntervalSeconds * time.Second):
				err := l.mutex.Renew(l.key)
				// TODO: Handle transient errors gently, don't just pass
				if err != nil {
					log.Debugf("Could not renew: %v", err)
					break
				}
				log.Debug("Lease renewed")
				break
			case wg = <-l.cancelRequired:
				log.Debug("Stopping renewing . . .")
				wg.Done()
				return
			}
		}
	}()
}

func (l *Locker) stopRenew() {
	var wg sync.WaitGroup
	wg.Add(1)
	l.cancelRequired <- &wg
	wg.Wait()
}
