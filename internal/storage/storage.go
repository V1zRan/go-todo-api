package storage

import (
	"encoding/json"
	"os"
	"todo-api/internal/models"
)

// Пока задачи хранятся в памяти.
// После подключения БД этот слайс будет удалён.
var tasks = []models.Task{
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
func CreateTask(title string) (models.Task, error) {
	task := models.Task{
		ID:    nextID,
		Title: title,
		Done:  false,
	}
	nextID++
	tasks = append(tasks, task)
	err := saveTasks()
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// Возвращает список всех задач.
func GetAllTasks() []models.Task {
	return tasks
}

// Ищет задачу по ID.
// Возвращает задачу и признак успешного поиска.
func GetTaskById(id int) (models.Task, bool) {
	for _, task := range tasks {
		if task.ID == id {
			return task, true
		}
	}
	return models.Task{}, false
}

// Обновляет задачу по ID.
// Возвращает обновлённую задачу и признак успешного обновления.
func UpdateTask(id int, req models.UpdateTaskRequest) (models.Task, bool, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = req.Title
			tasks[i].Done = req.Done

			err := saveTasks()
			if err != nil {
				return models.Task{}, true, err
			}

			return tasks[i], true, nil
		}
	}
	return models.Task{}, false, nil
}

// Удаляет задачу по ID.
// Возвращает true, если задача была найдена и удалена.
func DeleteTask(id int) (bool, error) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			err := saveTasks()
			if err != nil {
				return true, err
			}

			return true, nil
		}
	}
	return false, nil
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

func LoadTasks() error {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []models.Task{}
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
