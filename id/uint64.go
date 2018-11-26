package id

import (
	"fmt"
	"sync"
)

type innerUint64 struct {
	mu sync.Mutex
	id uint64
}

func (s *innerUint64) Max() uint64 {
	return maxUint64
}

func (s *innerUint64) New() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	if s.id > maxUint64 {
		s.id = 1
	}

	return s.id
}

func (s *innerUint64) NewAsString() string {
	return fmt.Sprint(s.New())
}
