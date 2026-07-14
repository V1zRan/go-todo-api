package postgres

import (
	"database/sql"
	"todo-api/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

// Функция подготавливает объект, через который в дальнейшем
// можно будет работать с PostgreSQL
func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) GetAllTasks() ([]models.Task, error) {
	// Выполняет SQL-запрос для получения всех задач
	rows, err := s.db.Query(`
	SELECT id, title, completed
	FROM tasks
	ORDER BY id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// Слайс для хранения задач, полученных из базы данных
	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		// Считывает текущую строку результата в структуру Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	// Проверяет наличие ошибок, возникших во время чтения результата запроса
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) CreateTask(title string) (models.Task, error) {
	var task models.Task
	// Передаём значения через параметры ($1), чтобы избежать SQL-инъекций.
	row := s.db.QueryRow(`
	INSERT INTO tasks (title)
	VALUES ($1)
	RETURNING id, title, completed;`, title)

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *Storage) GetTaskByID(id int) (models.Task, error) {
	var task models.Task

	row := s.db.QueryRow(`
	SELECT id, title, completed
	FROM tasks
	WHERE id = $1`, id)

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *Storage) UpdateTask(id int, req models.UpdateTaskRequest) (models.Task, error) {
	var task models.Task

	row := s.db.QueryRow(`
	UPDATE tasks
	SET title = $1,
		completed = $2
	WHERE id = $3
	RETURNING id, title, completed`, req.Title, req.Completed, id)

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *Storage) DeleteTask(id int) (bool, error) {
	result, err := s.db.Exec(`
	DELETE FROM tasks
	WHERE id = $1`, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}
