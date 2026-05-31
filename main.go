package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура описывает задачу, которая существует в системе
type Task struct {
	ID    int    `json:"id"`
	Title string `json"title"`
	Done  bool   `json"done"`
}

// Структура для создания задачи пользователем
type CreateTaskRequest struct {
	Title string `json:"title"`
}

// Пока задачи хранятся в памяти. после подключения БД
// После подключения БД этот слайс будет удалён.
var tasks = []Task{
	{
		ID:    1,
		Title: "30.05.2026 - изучать GO",
		Done:  false,
	},
	{
		ID:    2,
		Title: "30.05.2026 - Купить чоколодные трубочки",
		Done:  true,
	},
}

// Возвращает список задач в формате JSON
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	// Получение списка задач
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
		return
	}

	// Создание новой задачи
	if r.Method == http.MethodPost {
		// Структура для данных, пришедших от клиента
		var req CreateTaskRequest
		// Читаем JSON из тела запроса
		json.NewDecoder(r.Body).Decode(&req)
		// Создаем новую задачу
		newTask := Task{
			ID:    len(tasks) + 1,
			Title: req.Title,
			Done:  false,
		}
		tasks = append(tasks, newTask)
		// Возвращаем созданную задачу клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newTask)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {

	fmt.Println(tasks)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Home page")
		fmt.Fprintf(w, "Method: %s\n", r.Method)
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
	})

	http.HandleFunc("/tasks", tasksHandler)

	fmt.Println("Server os running on http://localhost:8080")

	// Запуск сервера
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}

}
