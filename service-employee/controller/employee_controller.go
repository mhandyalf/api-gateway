package controller

import (
	"fmt"
	"net/http"
	"service-employee/config"
	"service-employee/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var userURI string = "http://localhost:3001/user"

type WebResponse struct {
	Code   int
	Status string
	Data   interface{}
}

func CreateEmployee(c *fiber.Ctx) error {
	db := config.GetPostgresDB()

	var requestBody model.Employee

	c.BodyParser(&requestBody)

	requestBody.Id = uuid.New().String()

	accessToken := c.Get("access_token")
	if len(accessToken) == 0 {
		return c.Status(401).SendString("Invalid token: Access token missing")
	}

	req, err := http.NewRequest("GET", userURI+"/auth", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("access_token", accessToken)

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return c.Status(401).SendString("Invalid token")
	}

	ctx, cancel := config.NewPostgresContext()
	defer cancel()

	_, err = db.ExecContext(ctx, "INSERT INTO employee (id, name) VALUES ($1, $2)", requestBody.Id, requestBody.Name)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// Check if the error is a unique violation (duplicate key)
			if pgErr.Code == "23505" {
				return c.Status(400).SendString("Employee with the same ID already exists")
			}
		}
		return c.Status(500).SendString("Error inserting employee")
	}

	return c.JSON(WebResponse{
		Code:   201,
		Status: "OK",
		Data:   requestBody,
	})
}
