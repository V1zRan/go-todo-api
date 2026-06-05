// models.go
//
// Структуры данных приложения.
// Здесь описаны Task и структуры запросов API.
package main

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

// Пока задачи хранятся в памяти.
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
