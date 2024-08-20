package typeStruct

type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

func NewPost(title, content, link string, pubTime int64) Post {
	return Post{
		Title:   title,
		Content: content,
		PubTime: pubTime,
		Link:    link,
	}
}

type Pagination struct {
	Page       int
	PageSize   int
	TotalPages int
	TotalItems int
}

type PaginatedPosts struct {
	Posts      []Post
	Pagination Pagination
}
