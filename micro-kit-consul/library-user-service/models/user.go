package models

import "time"

type User struct {
	ID        uint64    `gorm:"primary_key" json:"id" form:"id"`
	CreatedAt time.Time `form:"created_at" json:"created_at"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at"`
	Username  string
	Password  string
	Email     string
}

func (User) TableName() string {
	return "user"
}
