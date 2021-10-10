package transport

import (
	"com/justin/micro/kit/pkg/configs"
	"com/justin/micro/kit/pkg/registers"
	"com/justin/micro/kit/pkg/tracers"
	"com/justin/micro/kit/pkg/utils"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	gozipkin "github.com/openzipkin/zipkin-go"
)

func NewHttpHandler(ctx context.Context, cfg *configs.AppConfig, tracer *gozipkin.Tracer, logger log.Logger) *gin.Engine {
	r := utils.NewRouter(ctx.Value("ginMod").(string))
	e := r.Group("/api/user")
	{
		e.POST("register", func(c *gin.Context) {
			register := registers.HttpClient(cfg, "user-service", "POST", "/api/v1/register", tracer, logger)
			kithttp.NewServer(
				register,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-user-servcie-register")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("findByID", func(c *gin.Context) {
			findByID := registers.HttpClient(cfg, "user-service", "GET", "/api/v1/findByID", tracer, logger)
			kithttp.NewServer(
				findByID,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-user-servcie-findByID")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("findByEmail", func(c *gin.Context) {
			findByEmail := registers.HttpClient(cfg, "user-service", "GET", "/api/v1/findByEmail", tracer, logger)
			kithttp.NewServer(
				findByEmail,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-user-servcie-findByEmail")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("findBooksByUserID", func(c *gin.Context) {
			findBooksByUserID := registers.HttpClient(cfg, "user-service", "GET", "/api/v1/findBooksByUserID", tracer, logger)
			kithttp.NewServer(
				findBooksByUserID,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-user-servcie-findBooksByUserID")...,
			).ServeHTTP(c.Writer, c.Request)
		})
	}

	e = r.Group("/api/book")
	{
		e.POST("save", func(c *gin.Context) {
			save := registers.HttpClient(cfg, "book-service", "POST", "/api/v1/save", tracer, logger)
			kithttp.NewServer(
				save,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-book-service-save")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("books", func(c *gin.Context) {
			books := registers.HttpClient(cfg, "book-service", "GET", "/api/v1/books", tracer, logger)
			kithttp.NewServer(
				books,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-book-service-books")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("selectBookByName", func(c *gin.Context) {
			selectBookByName := registers.HttpClient(cfg, "book-service", "GET", "/api/v1/selectBookByName", tracer, logger)
			kithttp.NewServer(
				selectBookByName,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-book-service-selectBookByName")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.POST("borrowBook", func(c *gin.Context) {
			borrowBook := registers.HttpClient(cfg, "book-service", "POST", "/api/v1/borrowBook", tracer, logger)
			kithttp.NewServer(
				borrowBook,
				utils.DecodeJSONRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "apigateway-book-service-borrowBook")...,
			).ServeHTTP(c.Writer, c.Request)
		})
	}

	return r
}
