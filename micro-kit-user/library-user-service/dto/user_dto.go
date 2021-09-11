package dto

type UserInfo struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUser struct {
	Username string
	Password string
	Email    string
}

type RegisterRequest struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Email    string `form:"email" json:"email" validate:"required,email"`
}

type FindByIDRequest struct {
	ID string `form:"id" json:"id" validate:"required"`
}

type FindByEmailRequest struct {
	Email string `form:"email" json:"email" validate:"required"`
}

type FindBooksByUserIDRequest struct {
	UserID string `form:"userid" json:"userid" validate:"required"`
}
