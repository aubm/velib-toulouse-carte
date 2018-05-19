package http

import (
	"context"
	stdhttp "net/http"

	"google.golang.org/appengine/urlfetch"
)

type ClientProvider interface {
	Provide(ctx context.Context, opts ClientOptions) *stdhttp.Client
}

type AppEngineClientProvider struct{}

func (p *AppEngineClientProvider) Provide(ctx context.Context, opts ClientOptions) *stdhttp.Client {
	return urlfetch.Client(ctx)
}

type ClientOptions struct {
}
