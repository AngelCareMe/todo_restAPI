package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"todo_api/db"
)

// Структура обработчика
type Handler struct {
	DB *pgxpool.Pool
}

// Инициализция обработчика
func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{DB: pool}
}

// Создаем обработчик на создание задачи
func (h *Handler) CreateTaskHandler(a *fiber.Ctx) error {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := a.BodyParser(&input); err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if input.Title == "" {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}

	task, err := db.CreateTask(h.DB, input.Title, input.Description)
	if err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return a.Status(fiber.StatusCreated).JSON(task)
}

// Создаем обработчик на обновление задачи
func (h *Handler) UpdateTaskHandler(a *fiber.Ctx) error {
	id, err := strconv.Atoi(a.Params("id"))
	if err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := a.BodyParser(&input); err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if input.Title == "" {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}

	task, err := db.UpdateTask(h.DB, id, input.Title, input.Description, input.Status)
	if err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return a.JSON(task)
}

// Создаем обработчик на удаление задачи
func (h *Handler) DeleteTaskHandler(a *fiber.Ctx) error {
	id, err := strconv.Atoi(a.Params("id"))
	if err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := db.DeleteTask(h.DB, id); err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return a.SendStatus(fiber.StatusNoContent)
}

// Создаем обработчик на получение всех задач
func (h *Handler) GetTasksHandler(a *fiber.Ctx) error {
	tasks, err := db.GetAllTasks(h.DB)
	if err != nil {
		return a.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return a.JSON(tasks)
}
