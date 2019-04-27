package container

import (
	"fmt"
	"sync"
)

type SafeSlice struct {
	values []interface{}

	lock *sync.RWMutex
}

func NewSafeSlice(size int) *SafeSlice {
	return &SafeSlice{
		lock:   &sync.RWMutex{},
		values: make([]interface{}, size),
	}
}

func (s *SafeSlice) Append(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.values = append(s.values, value)
}

func (s *SafeSlice) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.values)
}

func (s *SafeSlice) Range(f func(interface{}) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, value := range s.values {
		if !f(value) {
			break
		}
	}
}

func (s *SafeSlice) Get(index int) (interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if index >= len(s.values) {
		return nil, fmt.Errorf("SafeSlice: out of bound: %v", index)
	}
	return s.values[index], nil
}

func (s *SafeSlice) Set(index int, value interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if index >= len(s.values) {
		return fmt.Errorf("SafeSlice: out of bound: %v", index)
	}
	s.values[index] = value
	return nil
}

func (s *SafeSlice) Cover(values []interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.values = values
}
