package service

import (
	"com/justin/micro/kit/library-book-service/dto"
	"com/justin/micro/kit/pkg/monitors"
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

func (mw metricsMiddleware) SaveBook(ctx context.Context, bookname string) (*dto.BookInfo, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveBook"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.BookService.SaveBook(ctx, bookname)
}

func (mw metricsMiddleware) SelectBooks(ctx context.Context) ([]dto.BookInfo, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SelectBooks"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.BookService.SelectBooks(ctx)
}

func (mw metricsMiddleware) SelectBookByName(ctx context.Context, bookname string) (*dto.BookInfo, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SelectBookByName"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.BookService.SelectBookByName(ctx, bookname)
}

func (mw metricsMiddleware) BorrowBook(ctx context.Context, userID, bookID uint64) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "BorrowBook"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.BookService.BorrowBook(ctx, userID, bookID)
}

func (mw metricsMiddleware) HealthCheck() bool {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.BookService.HealthCheck()
}
