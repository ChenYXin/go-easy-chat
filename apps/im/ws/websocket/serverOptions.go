package websocket

import "time"

type ServerOptions func(opt *serverOptions)

type serverOptions struct {
	Authentication
	pattern string

	maxConnectIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) serverOptions {
	o := serverOptions{
		Authentication: new(authentication),
		maxConnectIdle: defaultMaxConnectionIdle,
		pattern:        "/ws",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOptions) {
		opt.Authentication = auth
	}
}

func WithServerPatten(pattern string) ServerOptions {
	return func(opt *serverOptions) {
		opt.pattern = pattern
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOptions) {
		if maxConnectionIdle > 0 {
			opt.maxConnectIdle = maxConnectionIdle
		}
	}
}
