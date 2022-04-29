package store

import (
	"sync"
)

type SyncSlice struct {
	Lock   sync.RWMutex
	values []any
}

type SyncSliceItem struct {
	Index int
	Value any
}

func NewSyncSlice() *SyncSlice {
	return &SyncSlice{
		Lock:   sync.RWMutex{},
		values: []any{},
	}
}

func (s *SyncSlice) Iterate() <-chan SyncSliceItem {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	ch := make(chan SyncSliceItem)

	go func() {
		for index, value := range s.values {
			ch <- SyncSliceItem{
				Index: index,
				Value: value,
			}
		}
		close(ch)
	}()

	return ch
}

func (s *SyncSlice) Len() int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	return len(s.values)
}

func (s *SyncSlice) Set(values []any) {
	s.Lock.Lock()
	s.values = values
	s.Lock.Unlock()
}

func (s *SyncSlice) Append(values ...any) {
	s.Lock.Lock()
	s.values = append(s.values, values...)
	s.Lock.Unlock()
}
