package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (srv *Server) HealthCheck(c *fiber.Ctx) error {
	logrus.WithField("test-key", "testing").WithField("test-key-2", "testing-2").Info("testing health route")
	c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Welcome to Golang, Fiber, and GORM",
	})
	return nil
}
