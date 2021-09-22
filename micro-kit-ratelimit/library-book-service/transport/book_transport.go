package transport

import (
	"com/justin/micro/kit/library-book-service/dto"
	"com/justin/micro/kit/library-book-service/endpoint"
	"com/justin/micro/kit/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(ctx context.Context, bookEndpoints *endpoint.BookEndpoints) *gin.Engine {
	r := utils.NewRouter(ctx.Value("ginMod").(string))

	e := r.Group("/api/v1")
	{
		e.POST("save", func(c *gin.Context) {
			kithttp.NewServer(
				bookEndpoints.SaveEndpoint,
				decodeBookRquest,
				utils.EncodeJsonResponse,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("books", func(c *gin.Context) {
			kithttp.NewServer(
				bookEndpoints.SelectBooksEndpoint,
				decodeBooksRequest,
				utils.EncodeJsonResponse,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.GET("selectBookByName", func(c *gin.Context) {
			kithttp.NewServer(
				bookEndpoints.SelectBookByNameEndpoint,
				decodeBookRquest,
				utils.EncodeJsonResponse,
			).ServeHTTP(c.Writer, c.Request)
		})

		e.POST("borrowBook", func(c *gin.Context) {
			kithttp.NewServer(
				bookEndpoints.BorrowBookEndpoint,
				decodeBorrowBookRequest,
				utils.EncodeJsonResponse,
			).ServeHTTP(c.Writer, c.Request)
		})
	}
	return r
}

func decodeBookRquest(_ context.Context, r *http.Request) (interface{}, error) {
	var req *dto.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.Bookname, nil
}

func decodeBooksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return "", nil
}

func decodeBorrowBookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req *dto.BorrowBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	userID, _ := strconv.ParseUint(req.UserID, 10, 64)
	bookID, _ := strconv.ParseUint(req.BookID, 10, 64)
	return dto.BorrowBook{
		UserID: userID,
		BookID: bookID,
	}, nil
}
