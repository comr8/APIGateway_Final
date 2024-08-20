package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoNews/pkg/api"
	"GoNews/pkg/database"
	"GoNews/pkg/parse"
)

func main() {
	// Чтение конфигурационного файла
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config file:", err)
	}
	defer configFile.Close()

	var config parse.Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("Failed to decode config file:", err)
	}

	// Инициализация базы данных
	db := database.InitDB()
	database.ExecuteSchemaSQL(db)
	defer db.Close()

	// Создание API сервера
	port := os.Getenv("API_PORT")
	apiPort := port
	go func() {
		err := api.StartAPI(apiPort, db)
		if err != nil {
			log.Fatal("Error starting API server:", err)
		}
	}()

	// Создание канала для завершения
	stopCh := make(chan struct{})

	// Обработка сигнала завершения
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		close(stopCh)
	}()

	// Запуск обхода RSS-лент
	for _, rssLink := range config.RSSLinks {
		go startParsingRoutine(rssLink, config.RequestPeriod, stopCh, db)
	}

	// Ожидание завершения
	<-stopCh
	fmt.Println("Application has been stopped.")
}

// startParsingRoutine запускает рутину для обхода и парсинга RSS-ленты.
func startParsingRoutine(url string, period int, stopCh <-chan struct{}, db *sql.DB) {
	for {
		select {
		case <-stopCh:
			return
		default:
			// Парсинг RSS-ленты
			posts, err := parse.ParseRSS(url)
			if err != nil {
				log.Println("Failed to parse RSS:", err)
				continue
			}

			// Сохранение постов в базу данных
			for _, post := range posts {
				_, err := database.SaveToDB(post)
				if err != nil {
					log.Println("Failed to save post to DB:", err)
				}
			}

			// Ожидание перед следующим обходом
			time.Sleep(time.Duration(period) * time.Minute)
		}
	}
}
