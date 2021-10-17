package service

import (
	"com/justin/micro/kit/library-book-service/dao"
	"com/justin/micro/kit/library-book-service/dto"
	"com/justin/micro/kit/library-book-service/models"
	"context"
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

var (
	ErrBookNotFound = errors.New(" This book is not exist! ")
	ErrBookExisted  = errors.New(" This book is existed! ")
	ErrBookSave     = errors.New(" Save book failed! ")
)

type BookService interface {
	SaveBook(ctx context.Context, bookname string) (*dto.BookInfo, error)
	SelectBooks(ctx context.Context) ([]dto.BookInfo, error)
	SelectBookByName(ctx context.Context, bookname string) (*dto.BookInfo, error)
	BorrowBook(ctx context.Context, userID, bookID uint64) error
	HealthCheck() bool
}

type BookServiceImpl struct {
	bookDao dao.BookDao
}

func NewBookServiceImpl(bookDao dao.BookDao) BookService {
	return &BookServiceImpl{
		bookDao: bookDao,
	}
}

func (b *BookServiceImpl) SaveBook(ctx context.Context, bookname string) (*dto.BookInfo, error) {
	book, err := b.bookDao.FindByName(bookname)
	if book != nil {
		log.Println("This book is already exist!")
		return &dto.BookInfo{}, ErrBookExisted
	}
	if err == gorm.ErrRecordNotFound || err == nil {
		newBook := &models.Book{Bookname: bookname}
		err = b.bookDao.Save(newBook)
		if err != nil {
			return nil, ErrBookSave
		}
		return &dto.BookInfo{
			ID:       newBook.ID,
			Bookname: newBook.Bookname,
		}, nil
	}
	return nil, err
}

func (b *BookServiceImpl) SelectBooks(ctx context.Context) ([]dto.BookInfo, error) {
	books, err := b.bookDao.FindAll()
	if err != nil {
		return nil, ErrBookNotFound
	}

	newBooks := new([]dto.BookInfo)
	for _, book := range books {
		*newBooks = append(*newBooks, dto.BookInfo{ID: book.ID, Bookname: book.Bookname})
	}
	return *newBooks, nil
}

func (b *BookServiceImpl) SelectBookByName(ctx context.Context, bookname string) (*dto.BookInfo, error) {
	book, err := b.bookDao.FindByName(bookname)
	if err != nil {
		return nil, ErrBookNotFound
	}
	return &dto.BookInfo{
		ID:       book.ID,
		Bookname: book.Bookname,
	}, nil
}

func (b *BookServiceImpl) BorrowBook(ctx context.Context, userID, bookID uint64) error {
	return b.bookDao.BorrowBook(userID, bookID)
}

func (b *BookServiceImpl) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(BookService) BookService
