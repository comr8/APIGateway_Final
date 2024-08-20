package main

import (
	"CommentsService/db"
	"CommentsService/pkg/comments"
	"CommentsService/pkg/types"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	datbase := db.InitDB()
	db.ExecuteSchemaSQL(datbase)
	defer db.CloseDB()

	router := gin.Default()
	port := os.Getenv("API_PORT")

	router.POST("/create-comment", func(c *gin.Context) {
		var request types.Request

		err := c.BindJSON(&request)
		uniqueID := request.UniqueID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusBadRequest, err)
			return
		}

		commentID, err := comments.AddComment(request.NewsID, request.CommentText, request.ParentCommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"commentID": commentID, "message": "Comment added successfully"})
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusOK)
	})

	router.POST("/del-comment", func(c *gin.Context) {
		var request types.Request

		err := c.BindJSON(&request)
		uniqueID := request.UniqueID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusBadRequest, err)
			return
		}

		err = comments.DeleteComment(request.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"commentID": request.ID, "message": "Comment deleted successfully"})
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusOK)
	})

	router.POST("/get-comment", func(c *gin.Context) {
		var request types.Request

		err := c.BindJSON(&request)
		uniqueID := request.UniqueID

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusBadRequest, err)
			return
		}

		comment, err := comments.GetComment(request.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusInternalServerError, err)
			return
		}
		fmt.Println(comment)
		c.JSON(http.StatusOK, comment)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusOK)

	})

	router.POST("/get-comments", func(c *gin.Context) {
		var request types.Request

		err := c.BindJSON(&request)
		uniqueID := request.UniqueID

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusBadRequest, err)
			return
		}

		comments, err := comments.GetCommentsByNewsID(request.NewsID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, comments)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusOK)
	})

	router.Run(":" + port)
}
