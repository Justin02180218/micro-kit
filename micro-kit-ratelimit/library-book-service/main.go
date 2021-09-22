package main

import (
	"com/justin/micro/kit/library-book-service/dao"
	"com/justin/micro/kit/library-book-service/endpoint"
	"com/justin/micro/kit/library-book-service/service"
	"com/justin/micro/kit/library-book-service/transport"
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/databases"
	"com/justin/micro/kit/pkg/ratelimits"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/juju/ratelimit"
)

var confFile = flag.String("f", "book.yaml", "Book config file")

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

	bucket := ratelimit.NewBucket(time.Second*time.Duration(configs.Conf.RatelimitConfig.FillInterval), int64(configs.Conf.RatelimitConfig.Capacity))
	ratelimit := ratelimits.NewTokenBucketLimiter(bucket)

	bookDao := dao.NewBookDaoImpl()
	bookService := service.NewBookServiceImpl(bookDao)

	bookEndpoints := &endpoint.BookEndpoints{
		SaveEndpoint:             ratelimit(endpoint.MakeSaveEndpoint(bookService)),
		SelectBooksEndpoint:      ratelimit(endpoint.MakeSelectBooksEndpoint(bookService)),
		SelectBookByNameEndpoint: ratelimit(endpoint.MakeSelectBookByNameEndpoint(bookService)),
		BorrowBookEndpoint:       ratelimit(endpoint.MakeBorrowBookEndpoint(bookService)),
	}

	ctx = context.WithValue(ctx, "ginMod", configs.Conf.ServerConfig.Mode)
	r := transport.NewHttpHandler(ctx, bookEndpoints)

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
