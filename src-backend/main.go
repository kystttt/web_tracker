package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/url"
	"os"
	"time"
)

var habits = map[string][]string{}

func loadFromDisk() {
	file, err := os.ReadFile("habits.json")
	if err == nil {
		json.Unmarshal(file, &habits)
	}
}

func saveToDisk() {
	data, _ := json.MarshalIndent(habits, "", "  ")
	os.WriteFile("habits.json", data, 0644)
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	loadFromDisk()

	app.Get("/api/habits", func(c *fiber.Ctx) error {
		return c.JSON(habits)
	})

	app.Post("/api/habits", func(c *fiber.Ctx) error {
		type req struct {
			Name string `json:"name"`
		}
		var r req
		if err := c.BodyParser(&r); err != nil {
			return c.Status(400).SendString("Invalid body")
		}
		if _, exists := habits[r.Name]; !exists {
			habits[r.Name] = []string{}
			saveToDisk()
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Post("/api/habits/:name", func(c *fiber.Ctx) error {
		encodedName := c.Params("name")
		name, err := url.QueryUnescape(encodedName)
		if err != nil {
			return c.Status(400).SendString("Invalid habit name")
		}

		today := time.Now().Format("2006-01-02")

		if _, exists := habits[name]; !exists {
			habits[name] = []string{}
		}

		for _, d := range habits[name] {
			if d == today {
				return c.SendStatus(fiber.StatusOK)
			}
		}

		habits[name] = append(habits[name], today)
		saveToDisk()
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Listen(":3001")
}
