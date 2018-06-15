package id

import (
	"sync"
	"time"
	"fmt"
)

type innerTime struct {
	id  uint64
	max uint64
	mu  sync.Mutex
}

func (s *innerTime) Max() uint64 {
	return s.max
}

func (s *innerTime) New() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	if s.id > s.max {
		now := time.Now()
		idStart := uint64(0)
		idStart += uint64(now.Year()%100) * 100000000000000
		idStart += uint64(now.Month()) * 1000000000000
		idStart += uint64(now.Day()) * 10000000000
		idStart += uint64(now.Hour()) * 100000000
		idStart += uint64(now.Minute()) * 1000000
		idStart += uint64(now.Second()) * 10000

		s.id = idStart + 1
		s.max = idStart + 9999
	}

	return s.id
}

func (s *innerTime) NewAsString() string {
	return fmt.Sprint(s.New())
}
