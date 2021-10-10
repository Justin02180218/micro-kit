package configs

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type ServerConfig struct {
	Port int    `json:"port" yaml:"port"`
	Mode string `json:"mode" yaml:"mode"`
	Name string `json:"name" yaml:"name"`
}

type MySQLConfig struct {
	Host     string `json:"host" yaml:"host"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Port     string `json:"port" yaml:"port"`
	Db       string `json:"db" yaml:"db"`
	Debug    bool   `json:"debug" yaml:"debug"`
}

type RatelimitConfig struct {
	FillInterval int `json:"fillInterval" yaml:"fillInterval"`
	Capacity     int `json:"capacity" yaml:"capacity"`
}

type ConsulConfig struct {
	Addr     string `json:"addr" yaml:"addr"`
	Interval string `json:"interval" yaml:"interval"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	Client   struct {
		RetryMax     int `json:"retrymax" yaml:"retrymax"`
		RetryTimeout int `json:"retrytimeout" yaml:"retrytimeout"`
	}
}

type GRPCConfig struct {
	RetryMax     int    `json:"retrymax" yaml:"retrymax"`
	RetryTimeout int    `json:"retrytimeout" yaml:"retrytimeout"`
	Name         string `json:"name" yaml:"name"`
}

type HystrixConfig struct {
	StreamPort string `json:"streamport" yaml:"streamport"`
}

type ZipkinConfig struct {
	Url         string `json:"url" yaml:"url"`
	ServiceName string `json:"service_name" yaml:"service_name"`
	Reporter    struct {
		Timeout       int `json:"timeout" yaml:"timeout"`
		BatchSize     int `json:"batch_size" yaml:"batch_size"`
		BatchInterval int `json:"batch_interval" yaml:"batch_interval"`
		MaxBacklog    int `json:"max_backlog" yaml:"max_backlog"`
	}
}

type PrometheusConfig struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Subsystem string `json:"subsystem" yaml:"subsystem"`
}

type AppConfig struct {
	*ServerConfig     `json:"server" yaml:"server"`
	*MySQLConfig      `json:"mysql" yaml:"mysql"`
	*RatelimitConfig  `json:"ratelimit" yaml:"ratelimit"`
	*ConsulConfig     `json:"consul" yaml:"consul"`
	*GRPCConfig       `json:"grpc" yaml:"grpc"`
	*HystrixConfig    `json:"hystrix" yaml:"hystrix"`
	*ZipkinConfig     `json:"zipkin" yaml:"zipkin"`
	*PrometheusConfig `json:"prometheus" yaml:"prometheus"`
}

var Conf = new(AppConfig)

func Init(file string) error {
	yamlData, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err = yaml.Unmarshal(yamlData, Conf); err != nil {
		return err
	}
	return nil
}
