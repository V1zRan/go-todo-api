package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// Структура для обновления данных задачи пользователем
type UpdateTaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
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

// Получаем ID задачи и возвращаем только её данные
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:

		for _, task := range tasks {
			if task.ID == id {
				json.NewEncoder(w).Encode(task)
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	case http.MethodPut:
		var req UpdateTaskRequest

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Title = req.Title
				tasks[i].Done = req.Done

				json.NewEncoder(w).Encode(tasks[i])
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	case http.MethodDelete:
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return

			}
		}

		http.Error(w, "task not found", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	test()

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
