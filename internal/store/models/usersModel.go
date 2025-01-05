package models

type UsersModel struct {
	ID        int64  `json:"id"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}
