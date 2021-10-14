package main

import (
	"com/justin/micro/kit/library-apigateway/transport"
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/tracers"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
)

var confFile = flag.String("f", "apigateway.yaml", "user config file")

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

	ctx := context.Background()

	reporter := tracers.NewZipkinReporter(configs.Conf.ZipkinConfig)
	defer reporter.Close()
	tracer := tracers.NewZipkinTracer(configs.Conf.ZipkinConfig.ServiceName, reporter)

	ctx = context.WithValue(ctx, "ginMod", configs.Conf.ServerConfig.Mode)
	r := transport.NewHttpHandler(ctx, configs.Conf, tracer, logger)

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
