package dto

type BookInfo struct {
	ID       uint64 `json:"id"`
	Bookname string `json:"bookname"`
}

type BookRequest struct {
	Bookname string `form:"bookname" json:"bookname" validate:"required"`
}

type BorrowBook struct {
	UserID uint64
	BookID uint64
}

type BorrowBookRequest struct {
	UserID string `form:"userid" json:"userid" validate:"required"`
	BookID string `form:"bookid" json:"bookid" validate:"required"`
}
