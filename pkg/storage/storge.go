package storage

import (
	"sync"
)

type Storage struct {
	mu     *sync.RWMutex
	groups map[string]*kindStorage
}

func New() *Storage {
	return &Storage{
		mu:     &sync.RWMutex{},
		groups: map[string]*kindStorage{},
	}
}

func (s *Storage) Store(state State) {
	s.mu.Lock()
	group, ok := s.groups[state.GetMeta().GetGroup()]
	if !ok {
		group = newKindStorage()
		s.groups[state.GetMeta().GetGroup()] = group
	}
	s.mu.Unlock()
	s.mu.RLock()
	defer s.mu.RUnlock()
	group.Store(state)
}

func (s *Storage) LoadGroup(group string) (kinds *kindStorage, loaded bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	kinds, ok := s.groups[group]
	return kinds, ok
}

func (s *Storage) LoadKind(group, kind string) (states *StateMap, loaded bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	kinds, ok := s.groups[group]
	if !ok {
		return nil, ok
	}
	return kinds.LoadKind(kind)
}

func (s *Storage) Load(group, kind, name string) (state State, loaded bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	kinds, ok := s.groups[group]
	if !ok {
		return nil, ok
	}
	return kinds.Load(kind, name)
}

func (s *Storage) DeleteGroup(group string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.groups, group)
}

func (s *Storage) DeleteKind(group, kind string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	kinds, ok := s.groups[group]
	if !ok {
		return
	}
	kinds.DeleteKind(kind)
}

func (s *Storage) Delete(group, kind, name string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	kinds, ok := s.groups[group]
	if !ok {
		return
	}
	kinds.Delete(kind, name)
}

func (s *Storage) RangeGroup(fn func(group string, kinds *kindStorage) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for group, kinds := range s.groups {
		if !fn(group, kinds) {
			break
		}
	}
}

func (s *Storage) RangeKind(fn func(group, kind string, states *StateMap) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for group, kinds := range s.groups {
		var breakFlag bool
		kinds.RangeKind(func(kind string, states *StateMap) bool {
			if !fn(group, kind, states) {
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

func (s *Storage) Range(fn func(group, kind, name string, state State) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for group, kinds := range s.groups {
		var breakFlag bool
		kinds.Range(func(kind, name string, state State) bool {
			if !fn(group, kind, name, state) {
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
