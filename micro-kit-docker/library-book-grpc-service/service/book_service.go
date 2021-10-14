package service

import (
	"com/justin/micro/kit/library-book-grpc-service/dao"
	"context"

	pbbook "com/justin/micro/kit/protos/book"
)

type BookService interface {
	FindBooksByUserID(ctx context.Context, req *pbbook.BooksByUserIDRequest) (*pbbook.BooksResponse, error)
}

type BookServiceImpl struct {
	bookDao dao.BookDao
}

func NewBookServiceImpl(bookDao dao.BookDao) BookService {
	return &BookServiceImpl{
		bookDao: bookDao,
	}
}

func (b *BookServiceImpl) FindBooksByUserID(ctx context.Context, req *pbbook.BooksByUserIDRequest) (*pbbook.BooksResponse, error) {
	books, err := b.bookDao.FindBooksByUserID(req.UserID)
	if err != nil {
		return &pbbook.BooksResponse{}, err
	}

	pbbooks := new([]*pbbook.BookInfo)
	for _, book := range books {
		*pbbooks = append(*pbbooks, &pbbook.BookInfo{
			Id:       book.ID,
			Bookname: book.Bookname,
		})
	}
	return &pbbook.BooksResponse{
		Books: *pbbooks,
	}, nil
}

type ServiceMiddleware func(BookService) BookService
