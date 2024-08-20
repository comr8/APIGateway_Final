package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"GoNews/pkg/database"

	"github.com/gorilla/mux"
)

// API представляет собой структуру для управления API.
type API struct {
	r        *mux.Router // Роутер для маршрутов API
	db       *sql.DB     // База данных
	rssLinks []string    /// Список ссылок на RSS-каналы
}

// NewAPI создает новый экземпляр API.
func NewAPI(db *sql.DB) *API {
	api := &API{
		r:  mux.NewRouter(), // Инициализация роутера
		db: db,              // Подключение к базе данных
	}

	api.endpoints() // Установка маршрутов API
	return api
}

// ServeHTTP позволяет API удовлетворять интерфейсу http.Handler.
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.r.ServeHTTP(w, r)
}

// GetRouter возвращает роутер API.
func (api *API) GetRouter() *mux.Router {
	return api.r
}

// posts обрабатывает запрос на получение последних новостей.
func (api *API) posts(w http.ResponseWriter, r *http.Request) {

	requestID := r.URL.Query().Get("uniqueID")
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		http.Error(w, "Неверное количество новостей", http.StatusBadRequest)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusBadRequest)
		return
	}

	posts, err := database.GetLatestPosts(n)
	if err != nil {
		http.Error(w, "Не удалось получить новости", http.StatusInternalServerError)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (api *API) Allposts(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("uniqueID")
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"]) // Количество новостей
	if err != nil {
		http.Error(w, "Invalid number of news", http.StatusBadRequest)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(vars["page"]) // Номер страницы
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusBadRequest)
		return
	}

	// Запрос новостей с учетом смещения
	posts, err := database.GetPosts(page, n)
	if err != nil {
		http.Error(w, "Failed to get news", http.StatusInternalServerError)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (api *API) searchPosts(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("uniqueID")
	vars := mux.Vars(r)
	keyword := vars["str"]

	if keyword == "" {
		http.Error(w, "Missing 'keyword' parameter in the request", http.StatusBadRequest)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusBadRequest)
		return
	}

	posts, err := database.SearchPostsByKeyword(keyword)
	if err != nil || len(posts) < 1 {
		http.Error(w, "Failed to search posts by keyword: "+keyword, http.StatusInternalServerError)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	// Если запрос завершился успешно, записываем успешный лог
	log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (api *API) getPost(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("uniqueID")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["n"])
	if err != nil {
		http.Error(w, "Failed id", http.StatusInternalServerError)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusBadRequest)
		return
	}

	posts, err := database.ReadFromDB(id)
	if err != nil {
		http.Error(w, "Не удалось получить новость", http.StatusInternalServerError)
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), requestID, r.RemoteAddr, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}

// webAppHandler обрабатывает запросы для веб-приложения.
func (api *API) webAppHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./webapp")).ServeHTTP(w, r)
}

// endpoints устанавливает маршруты API.
func (api *API) endpoints() {
	// Маршрут для получения n последних новостей
	api.r.HandleFunc("/news/{n:[0-9]+}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	// Маршрут для получения новостей с пагинацией
	api.r.HandleFunc("/news/{n:[0-9]+}/{page:[0-9]+}", api.Allposts).Methods(http.MethodGet, http.MethodOptions)
	// Маршрут для поиска по названию
	api.r.HandleFunc("/search/{str:[^/]+}", api.searchPosts).Methods(http.MethodGet)
	//Маршрут для получения новости по ID
	api.r.HandleFunc("/id/{n:[0-9]+}", api.getPost).Methods(http.MethodGet, http.MethodOptions)
	// Маршрут для обслуживания веб-приложения
	api.r.PathPrefix("/").HandlerFunc(api.webAppHandler).Methods(http.MethodGet)

}

// StartAPI запускает API на указанном порту.
func StartAPI(port string, db *sql.DB) error {
	api := NewAPI(db)
	return http.ListenAndServe(":"+port, api)
}
