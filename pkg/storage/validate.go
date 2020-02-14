package storage

import (
	"errors"
	"fmt"
)

var systemKinds = map[string]bool{}

const SystemGroup = "system"

func RegisterSystemKind(kind string) {
	systemKinds[kind] = true
}

func IsSystemKind(kind string) bool {
	return systemKinds[kind]
}

func Validate(state State) error {
	meta := state.GetMeta()
	if meta.GetGroup() == "" {
		return errors.New("group is empty")
	}

	if meta.GetKind() == "" {
		return errors.New("kind is empty")
	}

	if meta.GetName() == "" {
		return errors.New("name is empty")
	}

	if meta.GetGroup() == SystemGroup {
		if !IsSystemKind(meta.GetKind()) {
			return fmt.Errorf("kind \"%s\" is not in system group", meta.GetKind())
		}
	}

	return nil
}
