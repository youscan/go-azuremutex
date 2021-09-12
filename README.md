# Distributed Mutex on Azure Lease Blobs

[![CI](https://github.com/youscan/go-azuremutex/actions/workflows/ci.yaml/badge.svg)](https://github.com/youscan/azure-mutex/actions/workflows/ci.yaml)

This package implements distributed lock available for multiple processes. Possible use-cases include exclusive access to shared data, leader election, etc.

Azure Storage supports the feature that guarantees exclusive access to the blobs using basic lease functionality: it's possible to acquire a fixed or infinite duration lease, renew and then release it. Thus it's possible to implement a basic distributed mutex on top of blobs.

## Basic Abstractions

`AzureMutex`

A wrapper on top of Azure blobs API encapsulates API interactions and provides three primary functions: Acquire, Renew, Release.

`Locker`

Top-level abstractions that support background renewal of the acquired lease.

## Locker Usage

Import the module

```
import "github.com/youscan/go-azuremutex"
```

Configuration

```
var options = azmutex.MutexOptions{
	AccountName:   "mystorageaccount",
	AccountKey:    "******",
	ContainerName: "locks",
}
```

Use `Locker` with automated renewal

```
locker := azmutex.NewLocker(options, "migration")

err = locker.Lock()
...
// As soon the lease was acquired we can safely do exclusive job
RunDatabaseMigration()
...
err = locker.Unlock()
```

Use `AzureMutex` with lease that will expire in 60 seconds

```
mutex := azmutex.NewMutex(options)

err = mutex.Acquire("migration", 60)
...
// This code is only safely to run within 60 seconds window
RunDatabaseMigration()
...
err = mutex.Release("migration")
```

## Reference

- [Azure / Architecture / Cloud Design Patterns / Leader Election pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/leader-election)
- [Microsoft Azure Development Cookbook Second Edition: Leasing a blob and implementing distributed locks](https://www.oreilly.com/library/view/microsoft-azure-development/9781782170327/ch03s16.html)
