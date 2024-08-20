package comments_test

import (
	"CommentsService/db"
	"CommentsService/pkg/comments"
	"CommentsService/pkg/types"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestComment(t *testing.T) {

	db.InitDB()

	defer db.CloseDB()

	// Ваши параметры для добавления комментария
	newsID := 1
	text := "Новый комментарий"
	parentCommentID := 0

	commentID, err := comments.AddComment(newsID, text, parentCommentID)

	if err != nil {
		t.Fatalf("Ошибка при добавлении комментария: %v", err)
	} else {
		log.Printf("Комментарий с ID=%v добавлен", commentID)
	}

	// Теперь извлечем этот комментарий по commentID
	addedComment, err := comments.GetComment(commentID)

	if err != nil {
		t.Fatalf("Ошибка при извлечении комментария: %v", err)
	}

	// Теперь сравниваем извлеченный комментарий с ожидаемыми данными
	if addedComment.NewsID != newsID || addedComment.CommentText != text || addedComment.ParentCommentID != parentCommentID {
		t.Fatalf("Извлеченный комментарий не соответствует ожидаемым данными")
	}

	err = comments.DeleteComment(commentID)

	if err != nil {
		t.Fatalf("Ошибка при удалении комментария: %v", err)
	} else {
		log.Printf("Комментарий с ID=%v удален", commentID)
	}

	// Проверка, что комментарий был удален
	_, err = comments.GetComment(commentID)

	if err == nil {
		t.Fatalf("Ожидалась ошибка, так как комментарий должен быть удален")
	}

}

func TestGetCommentsByNewsID(t *testing.T) {
	db.InitDB()

	defer db.CloseDB()

	newsID := 1
	commentsToAdd := []types.Comment{
		{NewsID: newsID, CommentText: "Комментарий 1", ParentCommentID: 0},
		{NewsID: newsID, CommentText: "Комментарий 2", ParentCommentID: 0},
		{NewsID: newsID, CommentText: "Комментарий 3", ParentCommentID: 0},
	}

	var addedCommentIDs []int

	for _, comment := range commentsToAdd {
		commentID, err := comments.AddComment(comment.NewsID, comment.CommentText, comment.ParentCommentID)
		if err != nil {
			t.Fatalf("Ошибка при добавлении комментария: %v", err)
		}
		addedCommentIDs = append(addedCommentIDs, commentID)
	}

	commentsRetrieved, err := comments.GetCommentsByNewsID(newsID)

	if err != nil {
		t.Fatalf("Ошибка при извлечении комментариев: %v", err)
	}

	// Проверим, что количество извлеченных комментариев совпадает с добавленными
	if len(commentsRetrieved) != len(commentsToAdd) {
		t.Fatalf("Количество извлеченных комментариев не совпадает с ожидаемым")
	} else {
		log.Printf("Количество извлеченных комментариев совпадает с ожидаемым")
	}

	// Проверим, что тексты извлеченных комментариев совпадают с ожидаемыми
	for i, comment := range commentsRetrieved {
		if comment.CommentText != commentsToAdd[i].CommentText {
			t.Fatalf("Текст комментария не соответствует ожидаемому")
		} else {
			log.Printf("Текст комментария соответствует ожидаемому")
		}
	}

	// Удалим все добавленные комментарии
	for _, commentID := range addedCommentIDs {
		err := comments.DeleteComment(commentID)
		if err != nil {
			t.Fatalf("Ошибка при удалении комментария: %v", err)
		}
	}
}
