package environment

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type EnvironmentRequest struct {
	Name string `json:"name"`
	Type string `json:"type"` // e.g., "production", "staging"
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(router fiber.Router, handler *Handler) {
	group := router.Group("/environments")
	group.Post("/", handler.CreateEnvironment)
}

func (h *Handler) CreateEnvironment(c *fiber.Ctx) error {
	var req EnvironmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	state, err := h.service.StartProvisioning(c.Context(), req.Name, req.Type)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start provisioning"})
	}

	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Environment provisioning started",
		"id":      state.TransactionID,
		"status":  state.Status,
	})
}
