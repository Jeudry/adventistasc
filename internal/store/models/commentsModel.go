package models

type CommentsModel struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
