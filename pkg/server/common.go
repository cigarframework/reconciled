package server

import (
	"strings"
)

func joinKeys(keys ...string) string {
	return strings.Join(keys, "/")
}

func parseKey(key string) (group, kind, name string) {
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		return
	}
	group = parts[0]
	kind = parts[1]
	name = parts[2]
	return
}
