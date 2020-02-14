package grpcserver

type options struct {
	streamBufferSize int
}

type optionFunc func(options *options) *options

func WithStreamBufferSize(streamBufferSize int) optionFunc {
	return func(options *options) *options {
		options.streamBufferSize = streamBufferSize
		return options
	}
}
