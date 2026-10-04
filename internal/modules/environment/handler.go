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


// swagger:route POST /environments Environments createEnvironment
// 
// Start provisioning a new cloud environment
// 
// Asynchronously triggers a SAGA workflow to provision network, database, storage, compute, and secrets.
// 
// Responses:
// 202: environmentResponse
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

// swagger:parameters createEnvironment
type environmentRequestWrapper struct {
	// in:body
	Body EnvironmentRequest
}

// swagger:response environmentResponse
type environmentResponseWrapper struct {
	// in:body
	Body struct {
		Message string `json:"message"`
		ID      string `json:"id"`
		Status  string `json:"status"`
	}
}
