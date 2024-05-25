package handlers

import (
	"log/slog"
	"strings"

	"github.com/copito/data_quality/src/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (cc *Handlers) GetMetrics(c *fiber.Ctx) error {
	var results []model.Metric
	queryResults := cc.DB.Find(&results)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
		})
	}

	return c.Status(200).JSON(&results)
}

func (cc *Handlers) GetMetricByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(&fiber.Map{
			"error": "Invalid parameter",
		})
	}

	var result model.Metric
	queryResults := cc.DB.First(&result, id)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
		})
	}
	cc.Logger.Info(
		"queried metric",
		slog.String("metric_identifier", result.GetIdentifier(cc.DB)),
	)

	return c.Status(200).JSON(&result)
}

// CreateMetric
func (cc *Handlers) CreateMetric(c *fiber.Ctx) error {
	payload := model.Metric{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing data",
		})
	}

	// TODO: validation
	if strings.Trim(payload.TemplatedCalculation, " ") == "" {
		return c.Status(400).JSON(&fiber.Map{
			"message": "error validation templated_calculation",
		})
	}

	err := cc.DB.Transaction(func(tx *gorm.DB) error {
		queryResults := tx.Create(&payload).Error
		if queryResults != nil {
			return queryResults
		}

		// TODO: Generate CRON schedule
		// Create schedule

		return nil
	})
	// If rollback happened
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"message": "error validation templated_calculation",
			"context": err.Error(),
		})
	}

	return c.Status(200).JSON(&payload)
}
