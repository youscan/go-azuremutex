package azuremutex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getMutexOptionsForTests(t *testing.T) MutexOptions {
	return MutexOptions{
		ContainerName:      "locks",
		UseStorageEmulator: true,
		LogFunc: func(message string) {
			t.Logf(message)
		},
	}
}

func TestMutex(t *testing.T) {
	mutex := NewMutex(getMutexOptionsForTests(t))

	err := mutex.Acquire("mutex", 15)
	assert.NoError(t, err)

	err = mutex.Acquire("mutex", 15)
	assert.Error(t, err)

	err = mutex.Renew("mutex")
	assert.NoError(t, err)

	err = mutex.Release("mutex")
	assert.NoError(t, err)

	err = mutex.Acquire("mutex", 60)
	assert.NoError(t, err)

	_ = mutex.Release("mutex")
}
