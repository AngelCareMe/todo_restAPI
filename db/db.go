package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"os"
	"time"
	"todo_api/models"
)

// Инициализация БД с подключением через pgx
func InitDB() (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf(".env file not found")
	}
	
	requiredEnv := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "DB_NAME", "DB_SSLMODE"}
	for _, env := range requiredEnv {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", env)
		}
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	return pool, nil
}

// Создание таблицы задач (без миграции)
func CreateTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
    )
`)

	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	return nil
}

// Запись нового таска в БД
func CreateTask(pool *pgxpool.Pool, title, description string) (*models.Task, error) {
	var task models.Task
	err := pool.QueryRow(context.Background(), `
INSERT INTO tasks (title, description)
VALUES ($1, $2)
RETURNING id, title, description, status, created_at, updated_at
`, title, description).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Created_at,
		&task.Updated_at,
	)
	if err != nil {
		return nil, fmt.Errorf("create task error: %w", err)
	}
	return &task, nil
}

// Удаление таска из БД
func DeleteTask(pool *pgxpool.Pool, id int) error {
	delCom, err := pool.Exec(context.Background(), `
DELETE FROM tasks WHERE id = $1
`, id)
	if err != nil {
		return fmt.Errorf("delete task error: %w", err)
	}
	if delCom.RowsAffected() == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}

// Обновление задачи
func UpdateTask(pool *pgxpool.Pool, id int, title, description, status string) (*models.Task, error) {
	var task models.Task
	err := pool.QueryRow(context.Background(), `
UPDATE tasks
SET title = $1, description = $2, status = $3, updated_at = $4
WHERE id = $5
RETURNING id, title, description, status, created_at, updated_at
`, title, description, status, time.Now(), id).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Created_at,
		&task.Updated_at,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, fmt.Errorf("update task error: %w", err)
	}
	return &task, nil
}

// Получение всех задач из БД
func GetAllTasks(pool *pgxpool.Pool) ([]models.Task, error) {
	rows, err := pool.Query(context.Background(), `
SELECT id, title, description, status, created_at, updated_at
FROM tasks
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, fmt.Errorf("get all tasks error: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Created_at,
			&task.Updated_at,
		)
		if err != nil {
			return nil, fmt.Errorf("get all tasks error: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all tasks error: %w", err)
	}

	return tasks, nil
}
