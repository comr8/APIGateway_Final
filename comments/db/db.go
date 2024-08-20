package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var Host = os.Getenv("DB_HOST")
var Port = os.Getenv("DB_PORT")
var User = os.Getenv("DB_USER")
var Password = os.Getenv("DB_PASSWORD")
var DBName = os.Getenv("DB_NAME")

func InitDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Host, Port, User, Password, DBName)

	var db *sql.DB
	var err error

	maxRetries := 5                  // Максимальное количество попыток соединения
	retryInterval := 5 * time.Second // Интервал между попытками

	for retries := 0; retries < maxRetries; retries++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Failed to open database connection: %v", err)
			time.Sleep(retryInterval) // Подождите перед следующей попыткой
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database: %v", err)
			db.Close() // Закройте соединение перед следующей попыткой
			time.Sleep(retryInterval)
			continue
		}

		// Успешное соединение
		DB = db
		return DB
	}

	log.Printf("Exhausted all connection retries, giving up.")
	return nil
}

func ExecuteSchemaSQL(db *sql.DB) {

	// Чтение содержимого schema.sql
	schemaSQL, err := ioutil.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Выполнение SQL-запросов из schema.sql
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
