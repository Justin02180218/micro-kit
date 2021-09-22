package main

import (
	"com/justin/micro/kit/library-book-grpc-service/dao"
	"com/justin/micro/kit/library-book-grpc-service/endpoint"
	"com/justin/micro/kit/library-book-grpc-service/service"
	"com/justin/micro/kit/library-book-grpc-service/transport"
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/databases"
	"com/justin/micro/kit/pkg/ratelimits"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pbbook "com/justin/micro/kit/protos/book"

	"github.com/juju/ratelimit"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "book_rpc.yaml", "book rpc config file")
var quiteChan = make(chan error, 1)

func main() {
	flag.Parse()

	err := configs.Init(*configFile)
	if err != nil {
		panic(err)
	}
	err = databases.InitMySql(configs.Conf.MySQLConfig)
	if err != nil {
		fmt.Println("load mysql failed")
	}

	ctx := context.Background()

	bucket := ratelimit.NewBucket(time.Second*time.Duration(configs.Conf.RatelimitConfig.FillInterval), int64(configs.Conf.RatelimitConfig.Capacity))
	ratelimit := ratelimits.NewTokenBucketLimiter(bucket)

	bookDao := dao.NewBookDaoImpl()
	bookService := service.NewBookServiceImpl(bookDao)
	endpoints := endpoint.BookEndpoints{
		FindBooksByUserIDEndpoint: ratelimit(endpoint.NewFindBooksByUserIDEndpoint(bookService)),
	}

	go func() {
		handler := transport.NewBookServer(ctx, endpoints)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(configs.Conf.ServerConfig.Port)))
		if err != nil {
			fmt.Println("listen tcp err", err)
			quiteChan <- err
			return
		}
		gRPCServer := grpc.NewServer()
		pbbook.RegisterBookServer(gRPCServer, handler)
		err = gRPCServer.Serve(listener)
		if err != nil {
			fmt.Println("gRPCServer Serve err", err)
			quiteChan <- err
			return
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
		quiteChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-quiteChan)
}
