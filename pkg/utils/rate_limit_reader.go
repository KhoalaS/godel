package utils

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type RateLimitReader struct {
	Reader  io.Reader
	Ctx     context.Context
	Limiter *rate.Limiter
}

func (r *RateLimitReader) Read(p []byte) (int, error) {
	n := len(p)

	err := r.Limiter.WaitN(r.Ctx, n)
	if err != nil {
		return 0, err
	}

	return r.Reader.Read(p)
}
