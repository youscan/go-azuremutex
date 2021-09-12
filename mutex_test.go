package azuremutex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var options = MutexOptions{
	ContainerName:      "locks",
	UseStorageEmulator: true,
}

func TestMutex(t *testing.T) {
	mutex := NewMutex(options)

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
