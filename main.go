// main.go
//
// Точка входа приложения.
// Здесь настраиваются маршруты и запускается HTTP-сервер.
package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Home page")
		fmt.Fprintf(w, "Method: %s\n", r.Method)
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
	})

	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", handleTaskByID)

	fmt.Println("Server os running on http://localhost:8080")
	fmt.Println("page http://localhost:8080/tasks")

	// Запуск сервера
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}

}
