package monitors

import (
	"com/justin/micro/kit/pkg/configs"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type PrometheusParams struct {
	Counter   *kitprometheus.Counter
	Summary   *kitprometheus.Summary
	Gauge     *kitprometheus.Gauge
	Histogram *kitprometheus.Histogram
}

func MakePrometheusParams(conf *configs.PrometheusConfig) *PrometheusParams {
	fieldKeys := []string{"method"}
	requestCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: conf.Namespace,
		Subsystem: conf.Subsystem,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: conf.Namespace,
		Subsystem: conf.Subsystem,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	return &PrometheusParams{requestCounter, requestLatency, nil, nil}
}
