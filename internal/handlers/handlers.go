// handlers.go
//
// HTTP-обработчики для работы с задачами.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"todo-api/internal/models"
	"todo-api/internal/storage"
)

// tasksHandler обрабатывает запросы к списку задач.
// Поддерживает:
// GET /tasks  — получить все задачи
// POST /tasks — создать новую задачу
func TasksHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(storage.GetAllTasks())
		return
	}

	if r.Method == http.MethodPost {
		// Декодируем JSON из тела запроса в Go-структуру
		var req models.CreateTaskRequest

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

		newTask, err := storage.CreateTask(req.Title)
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
func HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	// Достаём ID задачи из URL и переводим его из строки в число.
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:

		task, found := storage.GetTaskById(id)
		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(task)
		return
	case http.MethodPut:
		var req models.UpdateTaskRequest

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

		updatedTask, found, err := storage.UpdateTask(id, req)
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
		deleted, err := storage.DeleteTask(id)

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
