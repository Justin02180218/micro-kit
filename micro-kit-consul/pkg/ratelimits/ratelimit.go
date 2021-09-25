package ratelimits

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/juju/ratelimit"
)

var ErrLimitExceed = errors.New(" Rate Limit Exceed ")

func NewTokenBucketLimiter(tb *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tb.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
