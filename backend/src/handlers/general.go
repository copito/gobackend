package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (cc *Handlers) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"success": true,
	})
}

func (cc *Handlers) Stack(c *fiber.Ctx) error {
	return c.JSON(c.App().Stack())
}
