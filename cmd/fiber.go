package cmd

import (
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/platform/cdi"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FiberApplication struct{}

func provideFiberApplication() *fiber.App {
	app := fiber.New()
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	
	dynamoClient := cdi.ProvideDynamoDB()
	sqsClient := cdi.ProvideSQS()

	tableName := "cloudforge-saga-state-local"
	queueUrl := "http://localhost:4566/000000000000/cloudforge-queue-local"

	repo := environment.NewRepository(dynamoClient, tableName)
	svc := environment.NewService(repo, sqsClient, queueUrl)
	handler := environment.NewHandler(svc)
	
	environment.RegisterRoutes(app, handler)

	return app
}

func (FiberApplication) Run() {
	app := provideFiberApplication()
	logrus.Info("Starting local Fiber server on port 8080")
	if err := app.Listen(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
