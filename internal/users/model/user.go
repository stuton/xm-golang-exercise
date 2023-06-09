package model

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
