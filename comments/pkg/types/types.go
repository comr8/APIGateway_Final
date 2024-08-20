package types

type Comment struct {
	ID              int    `json:"id"`
	NewsID          int    `json:"news_id"`
	CommentText     string `json:"commentText"`
	ParentCommentID int    `json:"parent_id"`
}

type Request struct {
	ID              int    `json:"id"`
	NewsID          int    `json:"news_id"`
	CommentText     string `json:"commentText"`
	ParentCommentID int    `json:"parent_id"`
	UniqueID        string `json:"uniqueID"`
}
