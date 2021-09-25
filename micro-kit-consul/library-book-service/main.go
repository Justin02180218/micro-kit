package main

import (
	"com/justin/micro/kit/library-book-service/dao"
	"com/justin/micro/kit/library-book-service/endpoint"
	"com/justin/micro/kit/library-book-service/service"
	"com/justin/micro/kit/library-book-service/transport"
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/databases"
	"com/justin/micro/kit/pkg/ratelimits"
	"com/justin/micro/kit/pkg/registers"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/juju/ratelimit"
)

var confFile = flag.String("f", "book.yaml", "Book config file")

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	err := configs.Init(*confFile)
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

	bookEndpoints := &endpoint.BookEndpoints{
		SaveEndpoint:             ratelimit(endpoint.MakeSaveEndpoint(bookService)),
		SelectBooksEndpoint:      ratelimit(endpoint.MakeSelectBooksEndpoint(bookService)),
		SelectBookByNameEndpoint: ratelimit(endpoint.MakeSelectBookByNameEndpoint(bookService)),
		BorrowBookEndpoint:       ratelimit(endpoint.MakeBorrowBookEndpoint(bookService)),
		HealthEndpoint:           ratelimit(endpoint.MakeHealthEndpoint(bookService)),
	}

	ctx = context.WithValue(ctx, "ginMod", configs.Conf.ServerConfig.Mode)
	r := transport.NewHttpHandler(ctx, bookEndpoints)

	check := registers.HttpCheck(configs.Conf.ServerConfig.Port, configs.Conf.ConsulConfig.Interval, configs.Conf.ConsulConfig.Timeout)
	registrar := registers.InitRegister(configs.Conf, check, logger)

	errChan := make(chan error)
	go func() {
		registrar.Register()
		errChan <- r.Run(fmt.Sprintf(":%s", strconv.Itoa(configs.Conf.ServerConfig.Port)))
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
	registrar.Deregister()
}
