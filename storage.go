package main

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
			return true
		}
	}
	return false
}
