package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s := New()

	state, ok := s.Load("group1", "kind1", "name1")
	assert.False(t, ok)
	assert.Nil(t, state)

	kinds, ok := s.LoadGroup("group1")
	assert.False(t, ok)
	assert.Nil(t, kinds)

	states, ok := s.LoadKind("group1", "kind1")
	assert.False(t, ok)
	assert.Nil(t, states)

	s.Store(NewMetaState("group1", "kind1", "name1"))
	state, ok = s.Load("group1", "kind1", "name1")
	assert.True(t, ok)
	assert.Equal(t, "group1", state.GetMeta().GetGroup())
	assert.Equal(t, "kind1", state.GetMeta().GetKind())
	assert.Equal(t, "name1", state.GetMeta().GetName())

	kinds, ok = s.LoadGroup("group1")
	assert.True(t, ok)
	state, ok = kinds.Load("kind1", "name1")
	assert.True(t, ok)
	assert.Equal(t, "group1", state.GetMeta().GetGroup())
	assert.Equal(t, "kind1", state.GetMeta().GetKind())
	assert.Equal(t, "name1", state.GetMeta().GetName())

	states, ok = s.LoadKind("group1", "kind1")
	assert.True(t, ok)
	state, ok = states.Load("name1")
	assert.True(t, ok)
	assert.Equal(t, "group1", state.GetMeta().GetGroup())
	assert.Equal(t, "kind1", state.GetMeta().GetKind())
	assert.Equal(t, "name1", state.GetMeta().GetName())

	s.Store(NewMetaState("group1", "kind1", "name2"))
	s.Store(NewMetaState("group1", "kind2", "name1"))
	s.Store(NewMetaState("group1", "kind2", "name2"))
	s.Store(NewMetaState("group2", "kind1", "name1"))
	s.Store(NewMetaState("group2", "kind1", "name2"))
	s.Store(NewMetaState("group2", "kind2", "name1"))
	s.Store(NewMetaState("group2", "kind2", "name2"))

	var list []State
	s.Range(func(_, _, _ string, state State) bool {
		list = append(list, state)
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("group1", "kind1", "name1"),
			NewMetaState("group1", "kind1", "name2"),
			NewMetaState("group1", "kind2", "name1"),
			NewMetaState("group1", "kind2", "name2"),
			NewMetaState("group2", "kind1", "name1"),
			NewMetaState("group2", "kind1", "name2"),
			NewMetaState("group2", "kind2", "name1"),
			NewMetaState("group2", "kind2", "name2"),
		},
		list)

	list = nil
	s.RangeKind(func(_, _ string, states *stateMap) bool {
		states.Range(func(_ string, state State) bool {
			list = append(list, state)
			return true
		})
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("group1", "kind1", "name1"),
			NewMetaState("group1", "kind1", "name2"),
			NewMetaState("group1", "kind2", "name1"),
			NewMetaState("group1", "kind2", "name2"),
			NewMetaState("group2", "kind1", "name1"),
			NewMetaState("group2", "kind1", "name2"),
			NewMetaState("group2", "kind2", "name1"),
			NewMetaState("group2", "kind2", "name2"),
		},
		list)

	list = nil
	s.RangeGroup(func(_ string, kinds *kindStorage) bool {
		kinds.Range(func(_, _ string, state State) bool {
			list = append(list, state)
			return true
		})
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("group1", "kind1", "name1"),
			NewMetaState("group1", "kind1", "name2"),
			NewMetaState("group1", "kind2", "name1"),
			NewMetaState("group1", "kind2", "name2"),
			NewMetaState("group2", "kind1", "name1"),
			NewMetaState("group2", "kind1", "name2"),
			NewMetaState("group2", "kind2", "name1"),
			NewMetaState("group2", "kind2", "name2"),
		},
		list)

	s.Delete("group1", "kind2", "name1")
	list = nil
	s.Range(func(_, _, _ string, state State) bool {
		list = append(list, state)
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("group1", "kind1", "name1"),
			NewMetaState("group1", "kind1", "name2"),
			NewMetaState("group1", "kind2", "name2"),
			NewMetaState("group2", "kind1", "name1"),
			NewMetaState("group2", "kind1", "name2"),
			NewMetaState("group2", "kind2", "name1"),
			NewMetaState("group2", "kind2", "name2"),
		},
		list)

	s.DeleteKind("group1", "kind1")
	list = nil
	s.Range(func(_, _, _ string, state State) bool {
		list = append(list, state)
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("group1", "kind2", "name2"),
			NewMetaState("group2", "kind1", "name1"),
			NewMetaState("group2", "kind1", "name2"),
			NewMetaState("group2", "kind2", "name1"),
			NewMetaState("group2", "kind2", "name2"),
		},
		list)

	s.DeleteGroup("group2")
	list = nil
	s.Range(func(_, _, _ string, state State) bool {
		list = append(list, state)
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("group1", "kind2", "name2"),
		},
		list)
}
