package api

import (
	"APIgateway/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"APIgateway/pkg/consts"
)

// AddComment отправляет комментарий в сервис комментариев.
func AddComment(newsId, parentCommentId int, commentText, uniqueID string) error {
	commentRequest := types.Request{
		NewsID:          newsId,
		ParentCommentID: parentCommentId,
		CommentText:     commentText,
		UniqueID:        uniqueID,
	}

	requestBody, err := json.Marshal(commentRequest)
	if err != nil {
		return err
	}

	response, err := http.Post(consts.CommentServiceURL+"/create-comment", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add comment. Status code: %d", response.StatusCode)
	}

	return nil
}

func DeleteComment(commentID int, uniqueID string) error {
	deleteRequest := types.Request{
		ID:       commentID,
		UniqueID: uniqueID,
	}

	requestBody, err := json.Marshal(deleteRequest)
	if err != nil {
		return err
	}

	response, err := http.Post(consts.CommentServiceURL+"/del-comment", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to delete comment. Status code: %d", response.StatusCode)
	}

	return nil

}

func GetComment(commentID int, uniqueID string) (types.Comment, error) {
	getCommentRequest := types.Request{
		ID:       commentID,
		UniqueID: uniqueID,
	}

	requestBody, err := json.Marshal(getCommentRequest)
	if err != nil {
		return types.Comment{}, err
	}

	response, err := http.Post(consts.CommentServiceURL+"/get-comment", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return types.Comment{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types.Comment{}, fmt.Errorf("Failed to get comment. Status code: %d", response.StatusCode)
	}

	var comments types.Comment
	if err := json.NewDecoder(response.Body).Decode(&comments); err != nil {
		return types.Comment{}, err
	}
	return comments, nil
}

func GetCommentsByNewsID(newsID int, uniqueID string) ([]types.Comment, error) {
	getCommentRequest := types.Request{
		NewsID:   newsID,
		UniqueID: uniqueID,
	}

	requestBody, err := json.Marshal(getCommentRequest)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(consts.CommentServiceURL+"/get-comments", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get comments. Status code: %d", response.StatusCode)
	}

	var comments []types.Comment
	if err := json.NewDecoder(response.Body).Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}
