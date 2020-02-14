package api

import (
	"context"

	"github.com/cigarframework/reconciled/pkg/storage"
)

type Action string

const (
	GetAction    = "get"
	ListAction   = "list"
	CreateAction = "create"
	UpdateAction = "update"
	PatchAction  = "patch"
	RemoveAction = "remove"
)

func (a Action) String() string {
	return string(a)
}

func (a Action) Valid() bool {
	if a == GetAction || a == ListAction || a == CreateAction || a == UpdateAction || a == PatchAction || a == RemoveAction {
		return true
	}
	return false
}

type ReviewRequest struct {
	Action     Action
	State      storage.State
	Patch      []*Patch
	User       *User
	Expression string
}

type Plugin interface {
	IsPlugin()
	Audit(context.Context, *ReviewRequest) error
	Admission(context.Context, *ReviewRequest) (*ReviewRequest, error)
	Start(context.Context) error
	Stop(context.Context) error
}
