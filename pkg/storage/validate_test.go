package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert.EqualError(t, Validate(NewMetaState("", "", "")), "group is empty")
	assert.EqualError(t, Validate(NewMetaState("group", "", "")), "kind is empty")
	assert.EqualError(t, Validate(NewMetaState("group", "kind", "")), "name is empty")
	assert.Nil(t, Validate(NewMetaState("group", "kind", "name")))

	assert.EqualError(t, Validate(NewMetaState("system", "kind", "test")), "kind \"kind\" is not in system group")
	assert.Nil(t, Validate(NewMetaState("system", "User", "name")))
}
