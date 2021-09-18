package endpoint

import (
	"com/justin/micro/kit/library-book-service/dto"
	"com/justin/micro/kit/library-book-service/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type BookEndpoints struct {
	SaveEndpoint             endpoint.Endpoint
	SelectBooksEndpoint      endpoint.Endpoint
	SelectBookByNameEndpoint endpoint.Endpoint
	BorrowBookEndpoint       endpoint.Endpoint
}

func MakeSaveEndpoint(svc service.BookService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bookname := request.(string)
		bookInfo, err := svc.SaveBook(ctx, bookname)
		if err != nil {
			return nil, err
		}
		return bookInfo, nil
	}
}

func MakeSelectBooksEndpoint(svc service.BookService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bookInfos, err := svc.SelectBooks(ctx)
		if err != nil {
			return nil, err
		}
		return bookInfos, nil
	}
}

func MakeSelectBookByNameEndpoint(svc service.BookService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bookname := request.(string)
		bookInfo, err := svc.SelectBookByName(ctx, bookname)
		if err != nil {
			return nil, err
		}
		return bookInfo, nil
	}
}

func MakeBorrowBookEndpoint(svc service.BookService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		borrowBook := request.(dto.BorrowBook)
		err = svc.BorrowBook(ctx, borrowBook.UserID, borrowBook.BookID)
		if err != nil {
			return nil, err
		}
		return "success", nil
	}
}
