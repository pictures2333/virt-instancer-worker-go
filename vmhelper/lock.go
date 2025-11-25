package vmhelper

import (
	"fmt"
	"sync"
)

var lockListLock = new(sync.Mutex)
var lockList = make(map[string]*sync.Mutex)

func vmlock(VMUUID string) (err error) {
	var (
		lock *sync.Mutex
		ok   bool
	)

	lockListLock.Lock()
	defer lockListLock.Unlock()

	// no lock -> create one
	if lock, ok = lockList[VMUUID]; !ok {
		lock = new(sync.Mutex)
		lockList[VMUUID] = lock
	}

	// lock
	if lock.TryLock() {
		return nil
	} else {
		return fmt.Errorf("Lock failed")
	}
}

func vmunlock(VMUUID string) (err error) {
	lockListLock.Lock()
	defer lockListLock.Unlock()

	if lock, ok := lockList[VMUUID]; ok {
		// unlock and delete
		lock.Unlock()
		delete(lockList, VMUUID)
		return nil
	} else {
		return fmt.Errorf("Lock not found")
	}
}
