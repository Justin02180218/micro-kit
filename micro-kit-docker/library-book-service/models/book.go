package models

import "time"

type Book struct {
	ID        uint64    `gorm:"primary_key" json:"id" form:"id"`
	CreatedAt time.Time `form:"created_at" json:"created_at"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at"`
	Bookname  string
}

func (Book) TableName() string {
	return "book"
}
