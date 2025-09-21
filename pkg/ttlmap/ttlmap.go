package ttlmap

import (
	"sync"
	"time"
)

type TTLMap struct {
	m             sync.Mutex
	stop          chan bool
	items         map[string]*ttlItem
	intervalCheck time.Duration
}

func New(intervalCheck time.Duration) *TTLMap {
	obj := TTLMap{
		m:             sync.Mutex{},
		stop:          make(chan bool),
		items:         make(map[string]*ttlItem),
		intervalCheck: intervalCheck,
	}

	go obj.start()
	return &obj
}

func (s *TTLMap) start() {
	for {
		select {
		case <-s.stop:
			goto exit
		case <-time.Tick(s.intervalCheck):
			s.CleanExpiredItems()
		}
	}
exit:
}

func (s *TTLMap) Stop() {
	s.stop <- true
}

func (s *TTLMap) Set(key string, value any, expiration *time.Duration) {
	s.m.Lock()
	defer s.m.Unlock()

	if value != nil {
		s.items[key] = newItem(value, expiration)
	}
}

func (s *TTLMap) Get(key string) any {
	s.m.Lock()
	defer s.m.Unlock()

	if val, ok := s.items[key]; ok {
		return val.value
	}

	return nil
}

func (s *TTLMap) Del(key string) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.items, key)
}

func (s *TTLMap) CleanExpiredItems() {
	s.m.Lock()
	defer s.m.Unlock()

	for key, item := range s.items {
		if item.Expired() {
			delete(s.items, key)
		}
	}
}

func (s *TTLMap) CleanAllItems() {
	s.m.Lock()
	defer s.m.Unlock()

	clear(s.items)
}
