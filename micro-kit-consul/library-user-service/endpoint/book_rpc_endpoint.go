package endpoint

import (
	"context"
	"fmt"

	pbbook "com/justin/micro/kit/protos/book"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type BookRPCEndpoints struct {
	FindBooksEndpoint endpoint.Endpoint
}

func MakeFindBooksEndpoint(instance string) endpoint.Endpoint {
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	findBooksEndpoint := grpctransport.NewClient(
		conn, "book.Book", "FindBooksByUserID",
		encodeGRPCFindBooksRequest,
		decodeGRPCFindBooksResponse,
		pbbook.BooksResponse{},
	).Endpoint()
	return findBooksEndpoint
}

func encodeGRPCFindBooksRequest(_ context.Context, r interface{}) (interface{}, error) {
	userID := r.(uint64)
	return &pbbook.BooksByUserIDRequest{
		UserID: userID,
	}, nil
}

func decodeGRPCFindBooksResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbbook.BooksResponse)
	return resp.Books, nil
}
