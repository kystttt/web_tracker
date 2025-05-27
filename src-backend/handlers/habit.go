package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kystttt/web_tracker/db"
	"time"
)

func getUserID(c *fiber.Ctx) int {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int(claims["user_id"].(float64))
}

func CreateHabit(c *fiber.Ctx) error {
	userID := getUserID(c)
	var body struct{ Name string }
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	_, err := db.Pool.Exec(c.Context(), "INSERT INTO habits (user_id, name) VALUES ($1, $2)", userID, body.Name)
	return err
}

func MarkHabitDone(c *fiber.Ctx) error {
	userID := getUserID(c)
	habitID := c.Params("id")
	date := time.Now().Format("2006-01-02")
	_, err := db.Pool.Exec(c.Context(), "INSERT INTO habit_marks (habit_id, date) SELECT $1, $2 WHERE EXISTS (SELECT 1 FROM habits WHERE id=$1 AND user_id=$3)", habitID, date, userID)
	return err
}

func GetUserHabits(c *fiber.Ctx) error {
	userID := getUserID(c)
	rows, err := db.Pool.Query(c.Context(), "SELECT id, name FROM habits WHERE user_id=$1", userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var habits []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		habits = append(habits, map[string]interface{}{"id": id, "name": name})
	}
	return c.JSON(habits)
}

func GetStats(c *fiber.Ctx) error {
	userID := getUserID(c)
	query := `
	SELECT h.id, h.name, COUNT(*) AS streak
	FROM habits h
	JOIN habit_marks m ON h.id = m.habit_id
	WHERE h.user_id=$1 AND m.date >= CURRENT_DATE - INTERVAL '20 days'
	GROUP BY h.id, h.name
	HAVING COUNT(*) = 21
	`
	rows, err := db.Pool.Query(c.Context(), query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var streak int
		rows.Scan(&id, &name, &streak)
		stats = append(stats, map[string]interface{}{"id": id, "name": name, "streak": streak})
	}
	return c.JSON(stats)
}
