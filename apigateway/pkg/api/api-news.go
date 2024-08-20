package api

import (
	"APIgateway/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"APIgateway/pkg/consts"
)

func GetLatestPosts(n int, uniqueID string) ([]types.NewsShortDetailed, error) {
	// Формируем URL для запроса

	url := consts.NewsServiceURL + "/news/" + strconv.Itoa(n) + "?uniqueID=" + uniqueID

	// Выполняем GET-запрос
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get posts. Status code: %d", response.StatusCode)
	}

	var posts []types.NewsShortDetailed
	if err := json.NewDecoder(response.Body).Decode(&posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllposts(n, page int, uniqueID string) (types.PaginatedPosts, error) {
	// Формируем URL для запроса

	url := consts.NewsServiceURL + "/news/" + strconv.Itoa(n) + "/" + strconv.Itoa(page) + "?uniqueID=" + uniqueID

	// Выполняем GET-запрос
	response, err := http.Get(url)
	if err != nil {
		return types.PaginatedPosts{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types.PaginatedPosts{}, fmt.Errorf("Failed to get posts. Status code: %d", response.StatusCode)
	}

	var posts types.PaginatedPosts
	if err := json.NewDecoder(response.Body).Decode(&posts); err != nil {
		return types.PaginatedPosts{}, err
	}

	return posts, nil
}

func SearchPosts(str, uniqueID string) ([]types.NewsShortDetailed, error) {
	// Формируем URL для запроса
	url := consts.NewsServiceURL + "/search/" + str + "?uniqueID=" + uniqueID

	// Выполняем GET-запрос
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get posts. Status code: %d", response.StatusCode)
	}

	var posts []types.NewsShortDetailed
	if err := json.NewDecoder(response.Body).Decode(&posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetNewsById(id int, uniqueID string) (types.NewsFullDetailed, error) {
	// Формируем URL для запроса
	url := consts.NewsServiceURL + "/id/" + strconv.Itoa(id) + "?uniqueID=" + uniqueID

	// Выполняем GET-запрос
	response, err := http.Get(url)
	if err != nil {
		return types.NewsFullDetailed{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return types.NewsFullDetailed{}, fmt.Errorf("Failed to get posts. Status code: %d", response.StatusCode)
	}

	var posts types.NewsFullDetailed
	if err := json.NewDecoder(response.Body).Decode(&posts); err != nil {
		return types.NewsFullDetailed{}, err
	}

	return posts, nil
}

func GetPostById(newsID int, uniqueID string) (types.ResultPost, error) {
	// канал для передачи результатов и ошибок
	resultCh := make(chan types.PostDetails, 2) // Буфер равен числу запросов

	// Создайте WaitGroup для синхронизации горутин
	var wg sync.WaitGroup

	// Запустите горутину для выполнения запроса к сервису news
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Выполните запрос к сервису news и обработайте результат
		newsResult, err := GetNewsById(newsID, uniqueID)
		resultCh <- types.PostDetails{News: newsResult, Err: err}
	}()

	// Запустите горутину для выполнения запроса к сервису comments
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Выполните запрос к сервису comments и обработайте результат
		commentsResult, err := GetCommentsByNewsID(newsID, uniqueID)
		resultCh <- types.PostDetails{Comments: commentsResult, Err: err}
	}()

	// Ожидайте завершения обоих горутин
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Обработайте результаты из канала и агрегируйте их
	var postDetails types.ResultPost
	for result := range resultCh {
		if result.Err != nil {
			return types.ResultPost{}, result.Err
		}

		if result.News.ID > 0 {
			postDetails.News = result.News
		}

		if result.Comments != nil {
			postDetails.Comments = result.Comments
		}
	}

	return postDetails, nil
}
