package ogmios

// Options available to ogmios client
type Options struct {
	endpoint string
	logger   Logger
	pipeline int
}

// Option to cardano client
type Option func(*Options)

// WithEndpoint allows ogmios endpoint to set; defaults to ws://127.0.0.1:1337
func WithEndpoint(endpoint string) Option {
	return func(opts *Options) {
		opts.endpoint = endpoint
	}
}

// WithLogger allows custom logger to be specified
func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.logger = logger
	}
}

// WithPipeline allows number of pipelined ogmios requests to be provided
func WithPipeline(n int) Option {
	return func(opts *Options) {
		opts.pipeline = n
	}
}

func buildOptions(opts ...Option) Options {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	if options.endpoint == "" {
		options.endpoint = "ws://127.0.0.1:1337"
	}
	if options.logger == nil {
		options.logger = DefaultLogger
	}
	if options.pipeline <= 0 {
		options.pipeline = 50
	}
	return options
}
