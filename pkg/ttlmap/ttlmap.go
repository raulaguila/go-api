package ttlmap

import (
	"sync"
	"time"
)

type TTLMap struct {
	m     sync.Mutex
	stop  chan bool
	items map[string]*ttlItem
}

func New(intervalCheck time.Duration) *TTLMap {
	obj := TTLMap{
		m:     sync.Mutex{},
		stop:  make(chan bool),
		items: make(map[string]*ttlItem),
	}

	go obj.loopCheckExpire(intervalCheck)
	return &obj
}

func (s *TTLMap) loopCheckExpire(interval time.Duration) {
	for {
		select {
		case <-s.stop:
			goto exit
		case <-time.Tick(interval):
			s.m.Lock()
			for key, item := range s.items {
				if item.Expired() {
					delete(s.items, key)
				}
			}
			s.m.Unlock()
		}
	}
exit:
}

func (s *TTLMap) Stop() {
	s.stop <- true
}

func (s *TTLMap) Set(key string, value any, expiration time.Duration) {
	s.m.Lock()
	defer s.m.Unlock()

	if value != nil {
		s.items[key] = newItem(value, time.Now().Add(expiration), true)
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

func (s *TTLMap) Clear() {
	s.m.Lock()
	defer s.m.Unlock()

	clear(s.items)
}
