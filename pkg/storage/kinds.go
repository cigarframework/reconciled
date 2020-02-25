package storage

import "sync"

type kindStorage struct {
	mu    *sync.RWMutex
	kinds map[string]*StateMap
}

func newKindStorage() *kindStorage {
	return &kindStorage{
		mu:    &sync.RWMutex{},
		kinds: map[string]*StateMap{},
	}
}

func (s *kindStorage) Store(state State) {
	s.mu.Lock()
	store, ok := s.kinds[state.GetMeta().GetKind()]
	if !ok {
		store = &StateMap{}
		s.kinds[state.GetMeta().GetKind()] = store
	}
	s.mu.Unlock()
	s.mu.RLock()
	defer s.mu.RUnlock()
	store.Store(state.GetMeta().GetName(), state)
}

func (s *kindStorage) Load(kind, name string) (state State, loaded bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	states, ok := s.kinds[kind]
	if !ok {
		return nil, ok
	}
	return states.Load(name)
}

func (s *kindStorage) LoadKind(kind string) (states *StateMap, loaded bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	states, ok := s.kinds[kind]
	return states, ok
}

func (s *kindStorage) Delete(kind, name string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	states, ok := s.kinds[kind]
	if !ok {
		return
	}
	states.Delete(name)
}

func (s *kindStorage) DeleteKind(kind string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.kinds, kind)
}

func (s *kindStorage) Range(fn func(kind, name string, state State) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for kind, states := range s.kinds {
		var breakFlag bool
		states.Range(func(name string, state State) bool {
			if !fn(kind, name, state) {
				breakFlag = true
				return false
			}
			return true
		})
		if breakFlag {
			break
		}
	}
}

func (s *kindStorage) RangeKind(fn func(kind string, states *StateMap) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for kind, states := range s.kinds {
		if !fn(kind, states) {
			break
		}
	}
}
