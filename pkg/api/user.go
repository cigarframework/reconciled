package api

import "github.com/cigarframework/reconciled/pkg/storage"

type User struct {
	Name  string
	Group string
	Token string
}

func init() {
	storage.RegisterSystemKind("User")
}
