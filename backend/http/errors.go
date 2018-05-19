package http

import "github.com/pkg/errors"

var (
	ErrUnexpectedStatusCode = errors.New("http: unexpected status code")
)
