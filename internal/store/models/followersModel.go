package models

type FollowersModel struct {
	UserId     int64  `json:"userId"`
	FollowerId int64  `json:"followerId"`
	CreatedAt  string `json:"createdAt"`
}
