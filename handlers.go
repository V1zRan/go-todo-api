// handlers.go
//
// HTTP-обработчики для работы с задачами.
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// tasksHandler обрабатывает запросы к списку задач.
// Поддерживает:
// GET /tasks  — получить все задачи
// POST /tasks — создать новую задачу
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
		return
	}

	if r.Method == http.MethodPost {
		// Декодируем JSON из тела запроса в Go-структуру
		var req CreateTaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

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

// handleTaskByID обрабатывает запросы к конкретной задаче.
// Поддерживает:
// GET /tasks/{id}    — получить задачу
// PUT /tasks/{id}    — обновить задачу
// DELETE /tasks/{id} — удалить задачу
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	// Достаём ID задачи из URL и переводим его из строки в число.
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
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
