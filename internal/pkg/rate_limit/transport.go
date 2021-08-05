package ratelimit

import (
	"net/http"

	"golang.org/x/time/rate"
)

type Transport struct {
	Delegate    http.RoundTripper
	Ratelimiter *rate.Limiter
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()
	err = t.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	next := t.Delegate
	if next == nil {
		next = http.DefaultTransport
	}
	return next.RoundTrip(req)
}
