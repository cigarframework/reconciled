package api

import "errors"

var (
	ErrBadData         = errors.New("malformed request")
	ErrNotExist        = errors.New("not exist")
	ErrExist           = errors.New("state exist")
	ErrUnAuthenticated = errors.New("not authenticated")
	ErrUnAuthorized    = errors.New("not authorized")
)
