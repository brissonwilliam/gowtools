package gowasync

import "sync"

type ErrorSyncGroup struct {
	mutex sync.Mutex
	errs  []error
}

func (eg *ErrorSyncGroup) Add(e error) {
	eg.mutex.Lock()
	eg.errs = append(eg.errs, e)
	eg.mutex.Unlock()
}

func (eg *ErrorSyncGroup) GetFirst() error {
	if len(eg.errs) > 0 {
		return eg.errs[0]
	}
	return nil
}
