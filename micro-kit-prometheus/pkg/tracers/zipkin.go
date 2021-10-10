package tracers

import (
	"com/justin/micro/kit/pkg/configs"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/kit/transport/http"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func NewZipkinReporter(zipkinCfg *configs.ZipkinConfig) reporter.Reporter {
	return zipkinhttp.NewReporter(
		zipkinCfg.Url,
		zipkinhttp.Timeout(time.Duration(zipkinCfg.Reporter.Timeout)*time.Second),
		zipkinhttp.BatchSize(zipkinCfg.Reporter.BatchSize),
		zipkinhttp.BatchInterval(time.Duration(zipkinCfg.Reporter.BatchInterval)*time.Second),
		zipkinhttp.MaxBacklog(zipkinCfg.Reporter.MaxBacklog),
	)
}

func NewZipkinTracer(serviceName string, reporter reporter.Reporter) *gozipkin.Tracer {
	zEP, _ := gozipkin.NewEndpoint(serviceName, "")

	tracerOptions := []gozipkin.TracerOption{
		gozipkin.WithLocalEndpoint(zEP),
	}

	tracer, err := gozipkin.NewTracer(reporter, tracerOptions...)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tracer
}

func Zipkin(tracer *gozipkin.Tracer, name string) endpoint.Middleware {
	return zipkin.TraceEndpoint(tracer, name)
}

func MakeHttpServerOptions(zipkinTracer *gozipkin.Tracer, name string) []http.ServerOption {
	zipkinServiceTrace := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name(name))
	options := []http.ServerOption{zipkinServiceTrace}
	return options
}

func MakeHttpClientOptions(zipkinTracer *gozipkin.Tracer, name string) []http.ClientOption {
	zipkinClientTrace := zipkin.HTTPClientTrace(zipkinTracer, zipkin.Name(name))
	options := []http.ClientOption{zipkinClientTrace}
	return options
}

func MakeGrpcServerOptions(zipkinTracer *gozipkin.Tracer, name string) []grpc.ServerOption {
	zipkinServiceTrace := zipkin.GRPCServerTrace(zipkinTracer, zipkin.Name(name))
	options := []grpc.ServerOption{zipkinServiceTrace}
	return options
}

func MakeGrpcClientOptions(zipkinTracer *gozipkin.Tracer, name string) []grpc.ClientOption {
	zipkinClientTrace := zipkin.GRPCClientTrace(zipkinTracer, zipkin.Name(name))
	options := []grpc.ClientOption{zipkinClientTrace}
	return options
}
