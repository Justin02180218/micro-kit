package endpoint

import (
	"com/justin/micro/kit/library-book-grpc-service/service"
	"context"

	"github.com/go-kit/kit/endpoint"
	pbbook "com/justin/micro/kit/protos/book"
)

type BookEndpoints struct {
	FindBooksByUserIDEndpoint endpoint.Endpoint
}

func NewFindBooksByUserIDEndpoint(bookService service.BookService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pbbook.BooksByUserIDRequest)
		res, err := bookService.FindBooksByUserID(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
