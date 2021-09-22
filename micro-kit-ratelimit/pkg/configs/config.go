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

type AppConfig struct {
	*ServerConfig    `json:"server" yaml:"server"`
	*MySQLConfig     `json:"mysql" yaml:"mysql"`
	*RatelimitConfig `json:"ratelimit" yaml:"ratelimit"`
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
