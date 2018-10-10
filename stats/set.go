package stats

import "sync"

type Set struct {
	sync.RWMutex
	name string
	m    map[Key]struct{}
}

func NewSet(name string) *Set {
	return &Set{
		name: name,
		m:    make(map[Key]struct{}),
	}
}

func (s *Set) Add(keys ...Key) *Set {
	s.Lock()
	defer s.Unlock()
	for i := 0; i < len(keys); i++ {
		s.m[keys[i]] = struct{}{}
	}
	return s
}

func (s *Set) Remove(keys ...Key) *Set {
	s.Lock()
	defer s.Unlock()
	for i := 0; i < len(keys); i++ {
		delete(s.m, keys[i])
	}
	return s

}

func (s *Set) Record(item *Item) error {
	var keys []Key
	for key, _ := range s.m {
		keys = append(keys, key)
	}
	if len(keys) > 0 {
		if err := Record(item, keys...); err != nil {
			return err
		}
	}
	return nil
}
