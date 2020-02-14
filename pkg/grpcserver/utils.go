package grpcserver

import (
	"errors"

	"github.com/cigarframework/reconciled/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, api.ErrBadData) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, api.ErrNotExist) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, api.ErrExist) {
		return status.Error(codes.AlreadyExists, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
