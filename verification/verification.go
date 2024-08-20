package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ArrowWords = []string{"java", "python", "ruby"}

func CreateComment(c *gin.Context) (string, string) {
	// Извлекаем текст комментария из тела запроса
	var request struct {
		CommentText string `json:"commentText"`
		UniqueID    string `json:"uniqueID"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return "", ""
	}

	// Приводим весь текст комментария к нижнему регистру
	commentText := strings.ToLower(request.CommentText)
	uniqueID := request.UniqueID

	return uniqueID, commentText
}

// код для проверки комментария
func verifyComment(commentText string, arrowWords []string) bool {
	// Создаем регулярное выражение
	re := regexp.MustCompile(strings.Join(arrowWords, "|"))

	// Ищем все вхождения запрещенных слов в тексте комментария
	matches := re.FindAllString(commentText, -1)
	fmt.Println(matches)

	// Если найдены совпадения (запрещенные слова), возвращаем false
	if len(matches) > 0 {
		return false
	}

	// В противном случае, комментарий прошел проверку
	return true
}

func answer(c *gin.Context, uniqueID string, verified bool) {
	if verified {
		// Если комментарий прошел проверку, отправляем статус 200 и uniqueID
		c.JSON(http.StatusOK, gin.H{"uniqueID": uniqueID, "message": "Comment verified", "error": ""})
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusOK)
	} else {
		// Если комментарий не прошел проверку, отправляем статус 400 и uniqueID
		c.JSON(http.StatusBadRequest, gin.H{"uniqueID": uniqueID, "error": "Comment verification failed"})
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusBadRequest)
	}
}

func main() {
	router := gin.Default()
	port := os.Getenv("API_PORT")

	// Маршрут для обработки POST-запросов
	router.POST("/verify", func(c *gin.Context) {
		uniqueID, commentText := CreateComment(c)
		verified := verifyComment(commentText, ArrowWords)
		answer(c, uniqueID, verified)
	})

	router.Run(":" + port)
}
