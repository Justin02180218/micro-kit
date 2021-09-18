package dao

import (
	"com/justin/micro/kit/library-book-service/models"
	"com/justin/micro/kit/pkg/databases"
)

type BookDao interface {
	Save(book *models.Book) error
	FindAll() ([]models.Book, error)
	FindByName(name string) (*models.Book, error)
	BorrowBook(userID, bookID uint64) error
}

type BookDaoImpl struct{}

func NewBookDaoImpl() BookDao {
	return &BookDaoImpl{}
}

func (b *BookDaoImpl) Save(book *models.Book) error {
	return databases.DB.Create(book).Error
}

func (b *BookDaoImpl) FindAll() ([]models.Book, error) {
	books := new([]models.Book)
	err := databases.DB.Find(books).Error
	if err != nil {
		return nil, err
	}
	return *books, nil
}

func (b *BookDaoImpl) FindByName(name string) (*models.Book, error) {
	book := &models.Book{}
	err := databases.DB.Where("bookname = ?", name).First(book).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (b *BookDaoImpl) BorrowBook(userID, bookID uint64) error {
	sql := "INSERT INTO user_book (user_id, book_id) VALUES(?, ?)"
	return databases.DB.Exec(sql, userID, bookID).Error
}
