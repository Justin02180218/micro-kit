package transport

import (
	"com/justin/micro/kit/library-user-service/dto"
	"com/justin/micro/kit/library-user-service/endpoint"
	"com/justin/micro/kit/pkg/tracers"
	"com/justin/micro/kit/pkg/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	kithttp "github.com/go-kit/kit/transport/http"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints, tracer *gozipkin.Tracer) *gin.Engine {
	r := utils.NewRouter(ctx.Value("ginMod").(string))

	e := r.Group("/api/v1")
	{
		e.POST("register", func(c *gin.Context) {
			kithttp.NewServer(
				endpoints.RegisterEndpoint,
				decodeRegisterRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "transport-register")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("findByID", func(c *gin.Context) {
			kithttp.NewServer(
				endpoints.FindByIDEndpoint,
				decodeFindByIDRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "transport-findByID")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("findByEmail", func(c *gin.Context) {
			kithttp.NewServer(
				endpoints.FindByEmailEndpoint,
				decodeFindByEmailRequest,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "transport-findByEmail")...,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("findBooksByUserID", func(c *gin.Context) {
			kithttp.NewServer(
				endpoints.FindBooksByUserIDEndpoint,
				decodeFindBooksByUserID,
				utils.EncodeJsonResponse,
				tracers.MakeHttpServerOptions(tracer, "transport-findBooks")...,
			).ServeHTTP(c.Writer, c.Request)
		})
	}

	r.GET("/health", func(c *gin.Context) {
		kithttp.NewServer(
			endpoints.HealthEndpoint,
			decodeHealthRequest,
			utils.EncodeJsonResponse,
			tracers.MakeHttpServerOptions(tracer, "transport-health")...,
		).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	return r
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req *dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &dto.RegisterUser{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}, nil
}

func decodeFindByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req *dto.FindByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.ID, nil
}

func decodeFindByEmailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req *dto.FindByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.Email, nil
}

func decodeFindBooksByUserID(_ context.Context, r *http.Request) (interface{}, error) {
	var req *dto.FindBooksByUserIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.UserID, nil
}

func decodeHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return struct{}{}, nil
}
