package websocket

type ServerOptions func(opt *serverOptions)

type serverOptions struct {
	Authentication
	patten string
}

func newServerOptions(opts ...ServerOptions) serverOptions {
	o := serverOptions{
		Authentication: new(authentication),
		patten:         "/ws",
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

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOptions) {
		opt.patten = patten
	}
}
