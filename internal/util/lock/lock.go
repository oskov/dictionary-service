package lock

import (
	"sync"
	"time"
)

type lock struct {
	mutex      Mutex
	accessTime time.Time
}

type LockStorage[K comparable] struct {
	mutexes map[K]lock
	rwMutex sync.RWMutex
	done    bool
}

func NewLockStorage[K comparable]() *LockStorage[K] {
	s := &LockStorage[K]{
		mutexes: make(map[K]lock),
	}

	go s.cleanUp()
	return s
}

func (s *LockStorage[K]) cleanUp() {
	for !s.done {
		time.Sleep(time.Minute)
		s.cleanUpMutexes()
	}
}

func (s *LockStorage[K]) cleanUpMutexes() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	for key, lock := range s.mutexes {
		if time.Since(lock.accessTime) > time.Minute {
			delete(s.mutexes, key)
		}
	}
}

func (s *LockStorage[K]) GetMutex(key K) Mutex {
	s.rwMutex.RLock()
	v, ok := s.mutexes[key]
	if ok {
		v.accessTime = time.Now()
		s.rwMutex.RUnlock()
		return v.mutex
	}

	s.rwMutex.RUnlock()
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.mutexes[key] = lock{
		mutex:      NewChanMutex(),
		accessTime: time.Now(),
	}

	return s.mutexes[key].mutex
}

func (s *LockStorage[K]) Close() {
	s.done = true
}

type Mutex interface {
	Lock()
	TryLock() bool
	TryLockWithTimeout(time.Duration) bool
	Unlock()
}

type ChanMutex struct {
	ch chan struct{}
}

func NewChanMutex() *ChanMutex {
	return &ChanMutex{
		ch: make(chan struct{}, 1),
	}
}

func (m *ChanMutex) Lock() {
	m.ch <- struct{}{}
}

func (m *ChanMutex) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (m *ChanMutex) TryLockWithTimeout(timeout time.Duration) bool {
	select {
	case m.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (m *ChanMutex) Unlock() {
	<-m.ch
}
