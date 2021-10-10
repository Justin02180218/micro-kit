package service

import (
	"com/justin/micro/kit/pkg/monitors"
	pbbook "com/justin/micro/kit/protos/book"
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

func MetricsMiddleware(prometheusParams *monitors.PrometheusParams) ServiceMiddleware {
	return func(next BookService) BookService {
		return metricsMiddleware{next, prometheusParams.Counter, prometheusParams.Summary}
	}
}

type metricsMiddleware struct {
	BookService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (mw metricsMiddleware) FindBooksByUserID(ctx context.Context, req *pbbook.BooksByUserIDRequest) (*pbbook.BooksResponse, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindBooksByUserID"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.BookService.FindBooksByUserID(ctx, req)
}
