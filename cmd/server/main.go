// main.go
//
// Точка входа приложения.
// Здесь настраиваются маршруты и запускается HTTP-сервер.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-api/internal/handlers"
	"todo-api/internal/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env для локальной разработки
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	postgresStorage, err := postgres.NewStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Home page")
		fmt.Fprintf(w, "Method: %s\n", r.Method)
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
	})

	http.HandleFunc("/tasks", handlers.TasksHandler(postgresStorage))
	http.HandleFunc("/tasks/", handlers.HandleTaskByID(postgresStorage))

	fmt.Println("Server os running on http://localhost:8080")
	fmt.Println("page http://localhost:8080/tasks")

	// Запуск сервера
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}

}
