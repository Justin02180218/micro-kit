package main

import (
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/databases"
	"flag"
	"fmt"
)

var confFile = flag.String("f", "user.yaml", "user config file")

func main() {
	flag.Parse()

	err := configs.Init(*confFile)
	if err != nil {
		panic(err)
	}

	err = databases.InitMySql(configs.Conf.MySQLConfig)
	if err != nil {
		fmt.Println("load mysql failed")
	}
	_, ok := databases.DB.Get("library")
	fmt.Println(ok)
}
