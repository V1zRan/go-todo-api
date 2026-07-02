// main.go
//
// Точка входа приложения.
// Здесь настраиваются маршруты и запускается HTTP-сервер.
package main

import (
	"fmt"
	"net/http"
	"todo-api/internal/handlers"
	"todo-api/internal/storage"
)

func main() {
	err := storage.LoadTasks()
	if err != nil {
		fmt.Println("Failed to load tasks", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Home page")
		fmt.Fprintf(w, "Method: %s\n", r.Method)
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
	})

	http.HandleFunc("/tasks", handlers.TasksHandler)
	http.HandleFunc("/tasks/", handlers.HandleTaskByID)

	fmt.Println("Server os running on http://localhost:8080")
	fmt.Println("page http://localhost:8080/tasks")

	// Запуск сервера
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}

}
