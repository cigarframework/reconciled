package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKinds(t *testing.T) {
	s := newKindStorage()

	state, ok := s.Load("kind1", "name1")
	assert.False(t, ok)
	assert.Nil(t, state)

	states, ok := s.LoadKind("kind1")
	assert.False(t, ok)
	assert.Nil(t, states)

	s.Store(NewMetaState("test", "kind1", "name1"))
	state, ok = s.Load("kind1", "name1")
	assert.True(t, ok)
	assert.Equal(t, "test", state.GetMeta().GetGroup())
	assert.Equal(t, "kind1", state.GetMeta().GetKind())
	assert.Equal(t, "name1", state.GetMeta().GetName())

	states, ok = s.LoadKind("kind1")
	assert.True(t, ok)
	state, ok = states.Load("name1")
	assert.True(t, ok)
	assert.Equal(t, "test", state.GetMeta().GetGroup())
	assert.Equal(t, "kind1", state.GetMeta().GetKind())
	assert.Equal(t, "name1", state.GetMeta().GetName())

	s.Store(NewMetaState("test", "kind1", "name2"))
	s.Store(NewMetaState("test", "kind2", "name1"))
	s.Store(NewMetaState("test", "kind2", "name2"))

	var list []State
	s.Range(func(_, _ string, state State) bool {
		list = append(list, state)
		return true
	})

	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("test", "kind1", "name1"),
			NewMetaState("test", "kind1", "name2"),
			NewMetaState("test", "kind2", "name1"),
			NewMetaState("test", "kind2", "name2"),
		},
		list)

	list = nil
	s.RangeKind(func(_ string, states *stateMap) bool {
		states.Range(func(_ string, state State) bool {
			list = append(list, state)
			return true
		})
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("test", "kind1", "name1"),
			NewMetaState("test", "kind1", "name2"),
			NewMetaState("test", "kind2", "name1"),
			NewMetaState("test", "kind2", "name2"),
		},
		list)

	s.Delete("kind1", "name2")
	state, ok = s.Load("kind1", "name2")
	assert.False(t, ok)
	assert.Nil(t, state)
	list = nil
	s.Range(func(_, _ string, state State) bool {
		list = append(list, state)
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("test", "kind1", "name1"),
			NewMetaState("test", "kind2", "name1"),
			NewMetaState("test", "kind2", "name2"),
		},
		list)

	s.DeleteKind("kind2")
	list = nil
	s.Range(func(_, _ string, state State) bool {
		list = append(list, state)
		return true
	})
	assert.ElementsMatch(
		t,
		[]State{
			NewMetaState("test", "kind1", "name1"),
		},
		list)
}
