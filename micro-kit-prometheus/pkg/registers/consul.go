package registers

import (
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/utils"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
	gozipkin "github.com/openzipkin/zipkin-go"
)

func InitRegister(cfg *configs.AppConfig, check api.AgentServiceCheck, logger log.Logger) (registrar sd.Registrar) {
	rand.Seed(time.Now().UnixNano())

	name := cfg.ServerConfig.Name
	addr := utils.LocalIP()
	port := cfg.ServerConfig.Port
	consulAddr := cfg.ConsulConfig.Addr
	client := connectConsul(consulAddr)
	num := rand.Intn(100)
	asr := api.AgentServiceRegistration{
		ID:      name + "-" + strconv.Itoa(num),
		Name:    name,
		Address: addr,
		Port:    port,
		Tags:    []string{name},
		Check:   &check,
	}
	registrar = consul.NewRegistrar(client, &asr, logger)
	return
}

func HttpCheck(port int, interval, timeout string) api.AgentServiceCheck {
	return api.AgentServiceCheck{
		HTTP:     "http://" + utils.LocalIP() + ":" + strconv.Itoa(port) + "/health",
		Interval: interval,
		Timeout:  timeout,
		Notes:    "Http Health Check",
	}
}

func GRPCCheck(port int, interval, timeout string) api.AgentServiceCheck {
	return api.AgentServiceCheck{
		GRPC:     utils.LocalIP() + ":" + strconv.Itoa(port) + "/health",
		Interval: interval,
		Timeout:  timeout,
		Notes:    "GRPC Health Check",
	}
}

func GRPCClient(cfg *configs.AppConfig, makeEndpoint func(string, *gozipkin.Tracer) endpoint.Endpoint, tracer *gozipkin.Tracer, logger log.Logger) endpoint.Endpoint {
	consulAddr := cfg.ConsulConfig.Addr
	retryMax := cfg.GRPCConfig.RetryMax
	retryTimeout := cfg.GRPCConfig.RetryTimeout
	name := cfg.GRPCConfig.Name

	client := connectConsul(consulAddr)
	instance := consul.NewInstancer(client, logger, name, []string{name}, true)
	factory := factoryFor(makeEndpoint, tracer)
	endpointer := sd.NewEndpointer(instance, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(retryMax, time.Millisecond*time.Duration(retryTimeout), balancer)
	return retry
}

func HttpClient(cfg *configs.AppConfig, name, method, path string, tracer *gozipkin.Tracer, logger log.Logger) endpoint.Endpoint {
	consulAddr := cfg.ConsulConfig.Addr
	retryMax := cfg.ConsulConfig.Client.RetryMax
	retryTimeout := cfg.ConsulConfig.Client.RetryTimeout

	client := connectConsul(consulAddr)
	instance := consul.NewInstancer(client, logger, name, []string{name}, true)
	factory := factoryForHttp(method, path, tracer)
	endpointer := sd.NewEndpointer(instance, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(retryMax, time.Millisecond*time.Duration(retryTimeout), balancer)
	return retry
}

func connectConsul(consulAddr string) (client consul.Client) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddr
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}
	client = consul.NewClient(consulClient)
	return
}

func factoryFor(makeEndpoint func(string, *gozipkin.Tracer) endpoint.Endpoint, tracer *gozipkin.Tracer) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		endpoint := makeEndpoint(instance, tracer)
		return endpoint, nil, nil
	}
}

func factoryForHttp(method, path string, tracer *gozipkin.Tracer) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		zipkinClientTrace := zipkin.HTTPClientTrace(tracer)
		options := []httptransport.ClientOption{zipkinClientTrace}
		return httptransport.NewClient(
			method,
			tgt,
			utils.EncodeJSONRequest,
			utils.DecodeJSONResponse,
			options...,
		).Endpoint(), nil, nil
	}
}
