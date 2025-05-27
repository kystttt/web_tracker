package models

type Habit struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type HabitMark struct {
	HabitID int    `json:"habit_id"`
	Date    string `json:"date"`
}
