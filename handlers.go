// handlers.go
//
// HTTP-обработчики для работы с задачами.
package main

import (
	"encoding/json"
	"errors"
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
		json.NewEncoder(w).Encode(getAllTasks())
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

		err = validateTitle(req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newTask, err := createTask(req.Title)
		if err != nil {
			// статус 500
			http.Error(w, "failed to save task", http.StatusInternalServerError)
			return
		}
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

		task, found := getTaskById(id)
		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(task)
		return
	case http.MethodPut:
		var req UpdateTaskRequest

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = validateTitle(req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedTask, found, err := updateTask(id, req)
		if err != nil {
			http.Error(w, "failed to save task", http.StatusInternalServerError)
		}

		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(updatedTask)
		return
	case http.MethodDelete:
		deleted, err := deleteTask(id)

		if err != nil {
			http.Error(w, "failed to delete task", http.StatusInternalServerError)
			return
		}

		if !deleted {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title is required")
	}

	if len([]rune(title)) > 100 {
		return errors.New("title is too long")
	}

	return nil
}
