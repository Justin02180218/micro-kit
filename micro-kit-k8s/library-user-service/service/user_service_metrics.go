package service

import (
	"com/justin/micro/kit/library-user-service/dto"
	"com/justin/micro/kit/pkg/monitors"
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

func MetricsMiddleware(prometheusParams *monitors.PrometheusParams) ServiceMiddleware {
	return func(next UserService) UserService {
		return metricsMiddleware{next, prometheusParams.Counter, prometheusParams.Summary}
	}
}

type metricsMiddleware struct {
	UserService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (mw metricsMiddleware) Register(ctx context.Context, vo *dto.RegisterUser) (*dto.UserInfo, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.UserService.Register(ctx, vo)
}

func (mw metricsMiddleware) FindByID(ctx context.Context, id uint64) (*dto.UserInfo, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindByID"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.UserService.FindByID(ctx, id)
}

func (mw metricsMiddleware) FindByEmail(ctx context.Context, email string) (*dto.UserInfo, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindByEmail"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.UserService.FindByEmail(ctx, email)
}

func (mw metricsMiddleware) FindBooksByUserID(ctx context.Context, id uint64) (interface{}, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindBooksByUserID"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.UserService.FindBooksByUserID(ctx, id)
}

func (mw metricsMiddleware) HealthCheck() bool {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.UserService.HealthCheck()
}
