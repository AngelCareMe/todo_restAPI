package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"todo_api/db"
	"todo_api/handlers"
)

func main() {
	//Инициализируем БД
	pool, err := db.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing DB: %s\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	//Создаем таблицу
	if err := db.CreateTable(pool); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating table: %s\n", err)
		os.Exit(1)
	}

	//Запускаем Fiber
	app := fiber.New()

	//Инициализируем ручки
	handler := handlers.NewHandler(pool)

	//Роут маршрутов
	app.Post("/tasks", handler.CreateTaskHandler)
	app.Get("/tasks", handler.GetTasksHandler)
	app.Put("/tasks/:id", handler.UpdateTaskHandler)
	app.Delete("/tasks/:id", handler.DeleteTaskHandler)

	//Запускаем сервер
	fmt.Println("Listening on port 8080")
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting app:", err)
	}
}
