package azmutex

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestConcurrentIncrement(t *testing.T) {
	const (
		threads    = 5
		operations = 10_000
	)

	var wg sync.WaitGroup
	wg.Add(threads)

	var count1 int
	var count2 int

	options := getMutexOptionsForTests(t)

	for i := 0; i < threads; i++ {
		go func() {
			locker := NewLocker(options, "test")
			err := locker.Lock()
			assert.NoError(t, err)
			if err != nil {
				return
			}

			for i := 0; i < operations; i++ {
				count1++
				count2 += 2
			}

			err = locker.Unlock()
			assert.NoError(t, err)

			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, operations*threads, count1)
	assert.Equal(t, operations*threads*2, count2)
}
