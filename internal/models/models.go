// models.go
//
// Структуры данных приложения.
// Здесь описаны Task и структуры запросов API.
package models

// Структура описывает задачу, которая существует в системе
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
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
