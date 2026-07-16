// main.go
//
// Точка входа приложения.
// Здесь настраиваются маршруты и запускается HTTP-сервер.
package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/handlers"
	"todo-api/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	postgresStorage, err := postgres.NewStorage(cfg.DatabaseURL)
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
