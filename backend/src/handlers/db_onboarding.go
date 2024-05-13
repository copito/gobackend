package handlers

import (
	"github.com/copito/data_quality/src/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (cc *Handlers) GetDatabaseOnboard(c *fiber.Ctx) error {
	var results []model.DatabaseOnboarding
	queryResults := cc.DB.Find(&results)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
		})
	}

	return c.Status(200).JSON(&results)
}

func (cc *Handlers) GetDatabaseOnboardByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(&fiber.Map{
			"error": "Invalid parameter",
		})
	}

	var result model.DatabaseOnboarding
	queryResults := cc.DB.First(&result, id)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
		})
	}

	return c.Status(200).JSON(&result)
}

// CreateDatabaseOnboard
func (cc *Handlers) CreateDatabaseOnboard(c *fiber.Ctx) error {
	payload := model.DatabaseOnboarding{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing data",
		})
	}

	// TODO: validation

	err := cc.DB.Transaction(func(tx *gorm.DB) error {
		queryResults := tx.Create(&payload).Error
		if queryResults != nil {
			return queryResults
		}

		// TODO: Test connection

		return nil
	})
	// If rollback happened
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"message": "error creating db_onboarding",
			"context": err.Error(),
		})
	}

	return c.Status(200).JSON(&payload)
}

func (cc *Handlers) CheckDatabaseOnboardConnectivity(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(&fiber.Map{
			"error": "Invalid parameter",
		})
	}

	var result model.DatabaseOnboarding
	queryResults := cc.DB.First(&result, id)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"message": "database could not be retrieved",
		})
	}

	// TODO: Check connectivity

	return c.Status(200).JSON(&result)
}
