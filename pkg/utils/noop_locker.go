package utils

type NoopLocker struct{}

func (l NoopLocker) Lock()    {}
func (l NoopLocker) Unlock()  {}
func (l NoopLocker) RLock()   {}
func (l NoopLocker) RUnlock() {}
