package environment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

type EnvironmentRequest struct {
	Name string `json:"name"`
	Type string `json:"type"` // e.g., "production", "staging"
}

func RegisterRoutes(router fiber.Router) {
	group := router.Group("/environments")
	group.Post("/", CreateEnvironment)
}

func CreateEnvironment(c *fiber.Ctx) error {
	var req EnvironmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// TODO: Call service to start SAGA workflow and save to DynamoDB
	// For now, mock successful reception
	envID := uuid.New().String()

	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Environment provisioning started",
		"id":      envID,
		"status":  "PENDING",
	})
}
