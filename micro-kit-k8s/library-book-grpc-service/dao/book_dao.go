package dao

import (
	"com/justin/micro/kit/library-book-grpc-service/models"
	"com/justin/micro/kit/pkg/databases"
)

type BookDao interface {
	FindBooksByUserID(userID uint64) ([]models.Book, error)
}

type BookDaoImpl struct{}

func NewBookDaoImpl() BookDao {
	return &BookDaoImpl{}
}

func (b *BookDaoImpl) FindBooksByUserID(userID uint64) ([]models.Book, error) {
	books := new([]models.Book)
	sql := "select b.* from book b, user_book ub where b.id = ub.book_id and ub.user_id = ?"
	err := databases.DB.Raw(sql, userID).Scan(books).Error
	if err != nil {
		return nil, err
	}
	return *books, nil
}
