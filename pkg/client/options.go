package client

import "google.golang.org/grpc"

type option struct {
	streamBufferSize int
	dialOptions      []grpc.DialOption
	getToken         func() string
}

type optionFunc func(optioptionons *option) *option

func WithStreamBufferSize(streamBufferSize int) optionFunc {
	return func(option *option) *option {
		option.streamBufferSize = streamBufferSize
		return option
	}
}

func WithDialoption(dialOptions ...grpc.DialOption) optionFunc {
	return func(option *option) *option {
		option.dialOptions = dialOptions
		return option
	}
}

func WithToken(getToken func() string) optionFunc {
	return func(option *option) *option {
		option.getToken = getToken
		return option
	}
}
