package main

import (
	"encoding/json"
	"os"
)

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

// Следующий ID для новой задачи.
var nextID = 3

// Создаёт новую задачу, сохраняет её в хранилище
// и возвращает созданный объект.
func createTask(title string) Task {
	task := Task{
		ID:    nextID,
		Title: title,
		Done:  false,
	}
	nextID++
	tasks = append(tasks, task)
	saveTasks()
	return task
}

// Возвращает список всех задач.
func getAllTasks() []Task {
	return tasks
}

// Ищет задачу по ID.
// Возвращает задачу и признак успешного поиска.
func getTaskById(id int) (Task, bool) {
	for _, task := range tasks {
		if task.ID == id {
			return task, true
		}
	}
	return Task{}, false
}

// Обновляет задачу по ID.
// Возвращает обновлённую задачу и признак успешного обновления.
func updateTask(id int, req UpdateTaskRequest) (Task, bool) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = req.Title
			tasks[i].Done = req.Done

			saveTasks()

			return tasks[i], true
		}
	}
	return Task{}, false
}

// Удаляет задачу по ID.
// Возвращает true, если задача была найдена и удалена.
func deleteTask(id int) bool {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			saveTasks()

			return true
		}
	}
	return false
}

const taskFile = "tasks.json"

func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(taskFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadTasks() error {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			nextID = 1
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &tasks)
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	nextID = maxID + 1
	if err != nil {
		return err
	}
	return nil
}
