package syncslice

import "sync"

type concurrentSliceItem[T comparable] struct {
	Index int
	Value T
}

type SyncSlice[T comparable] struct {
	mu    sync.RWMutex
	slice []T
}

func NewSyncSlice[T comparable](slice []T) *SyncSlice[T] {
	return &SyncSlice[T]{slice: slice}
}

func (s *SyncSlice[T]) Add(element T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice = append(s.slice, element)
}

func (s *SyncSlice[T]) Remove(element T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slice := []T{}

	for _, val := range s.slice {
		if element != val {
			slice = append(slice, val)
		}
	}

	s.slice = slice
}

func (s *SyncSlice[T]) Iter() <-chan concurrentSliceItem[T] {
	c := make(chan concurrentSliceItem[T])

	f := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for index, value := range s.slice {
			c <- concurrentSliceItem[T]{index, value}
		}
		close(c)
	}
	go f()

	return c
}
