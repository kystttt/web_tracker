package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kystttt/web_tracker/db"
	"github.com/kystttt/web_tracker/models"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"os"
	"strings"
	"time"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if strings.TrimSpace(user.Username) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username is required")
	}

	if _, err := mail.ParseAddress(user.Username); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}

	if len(user.Username) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Email too long")
	}
	if len(user.Password) < 6 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 6 characters")
	}
	var exists int
	err := db.Pool.QueryRow(c.Context(), "SELECT COUNT(*) FROM users WHERE username=$1", user.Username).Scan(&exists)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}
	if exists > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Password hashing failed")
	}
	user.Password = string(hash)

	_, err = db.Pool.Exec(c.Context(), "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	return c.SendStatus(fiber.StatusCreated)
}

func Login(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	var user models.User
	err := db.Pool.QueryRow(c.Context(), "SELECT id, password FROM users WHERE username=$1", input.Username).Scan(&user.ID, &user.Password)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		return fiber.ErrUnauthorized
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))

	return c.JSON(fiber.Map{"token": token})
}
