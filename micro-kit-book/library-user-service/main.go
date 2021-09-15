package main

import (
	"com/justin/micro/kit/library-user-service/dao"
	"com/justin/micro/kit/library-user-service/endpoint"
	"com/justin/micro/kit/library-user-service/service"
	"com/justin/micro/kit/library-user-service/transport"
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/databases"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	ctx := context.Background()

	userDao := dao.NewUserDaoImpl()
	userService := service.NewUserServiceImpl(userDao)
	userEndpoints := &endpoint.UserEndpoints{
		RegisterEndpoint:    endpoint.MakeRegisterEndpoint(userService),
		FindByIDEndpoint:    endpoint.MakeFindByIDEndpoint(userService),
		FindByEmailEndpoint: endpoint.MakeFindByEmailEndpoint(userService),
	}

	ctx = context.WithValue(ctx, "ginMod", configs.Conf.ServerConfig.Mode)
	r := transport.NewHttpHandler(ctx, userEndpoints)

	errChan := make(chan error)
	go func() {
		errChan <- r.Run(fmt.Sprintf(":%s", strconv.Itoa(configs.Conf.ServerConfig.Port)))
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}
