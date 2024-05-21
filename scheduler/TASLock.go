package scheduler

import (
	"sync/atomic"
)

type NewTASLock struct {
	lockStatus int64
}

func (myLock *NewTASLock) Lock() bool {
	return atomic.CompareAndSwapInt64(&myLock.lockStatus, 0, 1)
}

func (myLock *NewTASLock) Unlock() {
	atomic.StoreInt64(&myLock.lockStatus, 0)
}
