package stats

import (
	"context"
	"sync"

	ocstats "go.opencensus.io/stats"
)

type Set struct {
	sync.RWMutex
	name  string
	m     map[Key]context.Context
	cache []context.Context
}

func NewSet(name string) *Set {
	return &Set{
		name: name,
		m:    make(map[Key]context.Context),
	}
}
func (s *Set) updateCache() {
	s.RLock()
	s.RUnlock()
	var c []context.Context
	for _, ctx := range s.m {
		c = append(c, ctx)
	}
	s.cache = c
}
func (s *Set) Add(keys ...Key) *Set {
	s.Lock()
	defer s.Unlock()
	for i := 0; i < len(keys); i++ {
		s.m[keys[i]], _ = keys[i].context(context.Background())
	}
	var c []context.Context
	for _, ctx := range s.m {
		c = append(c, ctx)
	}
	s.cache = c
	return s
}

func (s *Set) Remove(keys ...Key) *Set {
	s.Lock()
	defer s.Unlock()
	for i := 0; i < len(keys); i++ {
		delete(s.m, keys[i])
	}
	var c []context.Context
	for _, ctx := range s.m {
		c = append(c, ctx)
	}
	s.cache = c
	return s

}

func (s *Set) Record(items ...Item) error {

	s.Lock()
	var localCache []context.Context
	localCache = append(localCache, s.cache...)
	s.Unlock()
	ms := getMeasurements(items...)
	go func(cache []context.Context, ms []ocstats.Measurement) {
		for _, ctx := range cache {
			ocstats.Record(ctx, ms...)
		}
	}(localCache, ms)
	return nil
}
