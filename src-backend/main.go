package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kystttt/web_tracker/db"
	"github.com/kystttt/web_tracker/handlers"
	"github.com/kystttt/web_tracker/middleware"
)

func main() {
	_ = godotenv.Load()
	db.Init()

	app := fiber.New()
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	api := app.Group("/api", middleware.JWTProtected())
	api.Post("/habit", handlers.CreateHabit)
	api.Post("/habit/:id/mark", handlers.MarkHabitDone)
	api.Get("/habits", handlers.GetUserHabits)
	api.Get("/stats", handlers.GetStats)

	app.Listen(":3000")
}
