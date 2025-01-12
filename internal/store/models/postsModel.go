package models

type PostsModel struct {
	ID        int64           `json:"id"`
	Content   string          `json:"content"`
	Title     string          `json:"title"`
	UserID    int64           `json:"user_id"`
	Tags      []string        `json:"tags"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Version   int             `json:"version"`
	Comments  []CommentsModel `json:"comments"`
	User      UsersModel      `json:"user"`
}

type PostWithMetadata struct {
	PostsModel
	CommentsCount int `json:"comment_count"`
}
