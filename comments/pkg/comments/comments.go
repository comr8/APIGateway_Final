package comments

import (
	"CommentsService/db"
	"CommentsService/pkg/types"
	"errors"
	"fmt"
)

// Метод для добавления комментария
func AddComment(newsID int, сommentText string, parentCommentID int) (int, error) {

	if parentCommentID != 0 && newsID != 0 {
		// Проверьте, что parentCommentID принадлежит указанной newsID
		query := "SELECT news_id FROM comments WHERE id = $1"
		var parentNewsID int
		err := db.DB.QueryRow(query, parentCommentID).Scan(&parentNewsID)
		if err != nil {
			return 0, err
		}

		if parentNewsID != newsID {
			return 0, errors.New("parentCommentID does not belong to the specified newsID")
		}
	}

	if newsID == 0 && parentCommentID != 0 {
		// Если newsID равно 0, выполните дополнительный запрос для извлечения newsID
		query := "SELECT news_id FROM comments WHERE id = $1"
		err := db.DB.QueryRow(query, parentCommentID).Scan(&newsID)
		if err != nil {
			return 0, err
		}
	}
	// SQL-запрос для вставки комментария в таблицу comments
	query := "INSERT INTO comments (news_id, text, parent_comment_id) VALUES ($1, $2, $3) RETURNING id"

	var commentID int
	err := db.DB.QueryRow(query, newsID, сommentText, parentCommentID).Scan(&commentID)

	if err != nil {
		return 0, err
	}

	return commentID, nil
}

// Метод для удаления комментария по ID
func DeleteComment(commentID int) error {
	// Проверим сначала, существует ли комментарий с указанным ID
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE id = $1", commentID).Scan(&count)
	if err != nil {
		return err
	}

	// Если комментарий с указанным ID не найден, вернем ошибку
	if count == 0 {
		return fmt.Errorf("comment with ID %d not found", commentID)
	}

	// Если комментарий с указанным ID существует, выполним удаление
	query := "DELETE FROM comments WHERE id = $1"
	_, err = db.DB.Exec(query, commentID)
	if err != nil {
		return err
	}

	return nil
}

// Метод извлечения комментария по ID
func GetComment(commentID int) (types.Comment, error) {
	var comment types.Comment
	err := db.DB.QueryRow("SELECT id, news_id, text, parent_comment_id FROM comments WHERE id = $1", commentID).Scan(&comment.ID, &comment.NewsID, &comment.CommentText, &comment.ParentCommentID)
	if err != nil {
		return types.Comment{}, err
	}
	return comment, nil
}

// Метод извлечения комментариев по ID новости
func GetCommentsByNewsID(newsID int) ([]types.Comment, error) {
	var comments []types.Comment
	rows, err := db.DB.Query("SELECT id, news_id, text, parent_comment_id FROM comments WHERE news_id = $1", newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment types.Comment
		if err := rows.Scan(&comment.ID, &comment.NewsID, &comment.CommentText, &comment.ParentCommentID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
