package types

type NewsShortDetailed struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NewsFullDetailed struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	PubTime int64  `json:"pubtime"`
	Link    string `json:"link"`
}

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
	Amount          int    `json:"Amount"`
}

type Pagination struct {
	Page       int
	PageSize   int
	TotalPages int
	TotalItems int
}

type PaginatedPosts struct {
	Posts      []NewsShortDetailed
	Pagination Pagination
}

type PostDetails struct {
	News     NewsFullDetailed
	Comments []Comment
	Err      error
}

type ResultPost struct {
	News     NewsFullDetailed
	Comments []Comment
}
