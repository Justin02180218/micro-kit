package registers

import (
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/utils"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hashicorp/consul/api"
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

func GRPCClient(cfg *configs.AppConfig, makeEndpoint func(string) endpoint.Endpoint, logger log.Logger) endpoint.Endpoint {
	consulAddr := cfg.ConsulConfig.Addr
	retryMax := cfg.GRPCConfig.RetryMax
	retryTimeout := cfg.GRPCConfig.RetryTimeout
	name := cfg.GRPCConfig.Name

	client := connectConsul(consulAddr)
	instance := consul.NewInstancer(client, logger, name, []string{name}, true)
	factory := factoryFor(makeEndpoint)
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

func factoryFor(makeEndpoint func(string) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		endpoint := makeEndpoint(instance)
		return endpoint, nil, nil
	}
}
