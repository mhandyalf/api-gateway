package controller

import (
	"context"
	"errors"
	"service-user/helpers"
	"service-user/model"
	"service-user/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WebResponse struct {
	Code   int
	Status string
	Data   interface{}
}

func Register(c *fiber.Ctx) error {
	var requestBody model.User
	db := config.GetPostgresDB()

	requestBody.Id = uuid.New().String()

	ctx, cancel := config.NewPostgresContext()
	defer cancel()

	c.BodyParser(&requestBody)

	hashedPassword := helpers.HashPassword([]byte(requestBody.Password))

	_, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", requestBody.Id, requestBody.Email, hashedPassword)
	if err != nil {
		return c.JSON(WebResponse{
			Code:   500,
			Status: "INTERNAL_SERVER_ERROR",
			Data:   err.Error(),
		})
	}

	return c.JSON(WebResponse{
		Code:   201,
		Status: "OK",
		Data:   requestBody.Email,
	})
}

func Login(c *fiber.Ctx) error {
	db := config.GetPostgresDB()

	var requestBody model.User
	var result model.User

	c.BodyParser(&requestBody)

	err := db.QueryRowContext(context.TODO(), "SELECT id, email, password FROM users WHERE email = $1", requestBody.Email).Scan(&result.Id, &result.Email, &result.Password)
	if err != nil {
		return c.JSON(WebResponse{
			Code:   401,
			Status: "UNAUTHORIZED",
			Data:   errors.New("invalid email").Error(),
		})
	}

	checkPassword := helpers.ComparePassword([]byte(result.Password), []byte(requestBody.Password))
	if !checkPassword {
		return c.JSON(WebResponse{
			Code:   401,
			Status: "UNAUTHORIZED",
			Data:   errors.New("invalid password").Error(),
		})
	}

	access_token := helpers.SignToken(requestBody.Email)

	return c.JSON(struct {
		Code        int
		Status      string
		AccessToken string
		Data        interface{}
	}{
		Code:        200,
		Status:      "OK",
		AccessToken: access_token,
		Data:        result,
	})
}

func Auth(c *fiber.Ctx) error {
	return c.JSON("OK")
}
