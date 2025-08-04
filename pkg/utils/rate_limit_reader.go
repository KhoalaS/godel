package utils

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type RateLimitReader struct {
	reader  io.Reader
	ctx     context.Context
	limiter *rate.Limiter
}

func (r *RateLimitReader) Read(p []byte) (int, error) {
	n := len(p)

	err := r.limiter.WaitN(r.ctx, n)
	if err != nil {
		return 0, err
	}

	return r.reader.Read(p)
}
