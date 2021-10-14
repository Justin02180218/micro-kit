package main

import (
	"com/justin/micro/kit/library-book-grpc-service/dao"
	"com/justin/micro/kit/library-book-grpc-service/endpoint"
	"com/justin/micro/kit/library-book-grpc-service/service"
	"com/justin/micro/kit/library-book-grpc-service/transport"
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/databases"
	"com/justin/micro/kit/pkg/monitors"
	"com/justin/micro/kit/pkg/ratelimits"
	"com/justin/micro/kit/pkg/registers"
	"com/justin/micro/kit/pkg/tracers"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pbbook "com/justin/micro/kit/protos/book"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"github.com/juju/ratelimit"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "book_rpc.yaml", "book rpc config file")
var quiteChan = make(chan error, 1)

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

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

	reporter := tracers.NewZipkinReporter(configs.Conf.ZipkinConfig)
	defer reporter.Close()
	tracer := tracers.NewZipkinTracer(configs.Conf.ZipkinConfig.ServiceName, reporter)

	prometheusParams := monitors.MakePrometheusParams(configs.Conf.PrometheusConfig)

	bookDao := dao.NewBookDaoImpl()
	bookService := service.NewBookServiceImpl(bookDao)
	bookService = service.MetricsMiddleware(prometheusParams)(bookService)
	endpoints := endpoint.BookEndpoints{
		FindBooksByUserIDEndpoint: tracers.Zipkin(tracer, "grpc-endpoint-findBooks")(ratelimit(endpoint.NewFindBooksByUserIDEndpoint(bookService))),
	}

	// check := registers.GRPCCheck(configs.Conf.ServerConfig.Port, configs.Conf.ConsulConfig.Interval, configs.Conf.ConsulConfig.Timeout)
	var check api.AgentServiceCheck
	registrar := registers.InitRegister(configs.Conf, check, logger)

	go func() {
		registrar.Register()
		handler := transport.NewBookServer(ctx, endpoints, tracer)
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
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/metrics",
			promhttp.Handler().ServeHTTP,
		)
		quiteChan <- http.ListenAndServe(":10089", mux)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
		quiteChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-quiteChan)
	registrar.Deregister()
}
