package transport

import (
	"com/justin/micro/kit/library-book-grpc-service/endpoint"
	"context"

	pbbook "com/justin/micro/kit/protos/book"

	kitrpc "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pbbook.UnimplementedBookServer
	findBooksByUserID kitrpc.Handler
}

func (g grpcServer) FindBooksByUserID(ctx context.Context, r *pbbook.BooksByUserIDRequest) (*pbbook.BooksResponse, error) {
	_, res, err := g.findBooksByUserID.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return res.(*pbbook.BooksResponse), nil
}

func NewBookServer(ctx context.Context, endpoints endpoint.BookEndpoints) pbbook.BookServer {
	return &grpcServer{
		findBooksByUserID: kitrpc.NewServer(
			endpoints.FindBooksByUserIDEndpoint,
			decodeFindBooksByUserIDRequest,
			encodeFindBooksByUserIDResponse,
		),
	}
}

func decodeFindBooksByUserIDRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*pbbook.BooksByUserIDRequest)
	return &pbbook.BooksByUserIDRequest{
		UserID: req.UserID,
	}, nil
}

func encodeFindBooksByUserIDResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pbbook.BooksResponse)
	return &pbbook.BooksResponse{
		Books: resp.Books,
	}, nil
}
