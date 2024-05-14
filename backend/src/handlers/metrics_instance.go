package handlers

import (
	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/adhocore/gronx"
)

func (cc *Handlers) GetMetricInstances(c *fiber.Ctx) error {
	var results []model.MetricInstance
	queryResults := cc.DB.Find(&results)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
		})
	}

	return c.Status(200).JSON(&results)
}

func (cc *Handlers) GetMetricInstanceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(&fiber.Map{
			"error": "Invalid parameter",
		})
	}

	var result model.MetricInstance
	queryResults := cc.DB.First(&result, id)

	if queryResults.Error != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
		})
	}

	return c.Status(200).JSON(&result)
}

// CreateMetricInstanceByID
func (cc *Handlers) CreateMetricInstanceByID(c *fiber.Ctx) error {
	payload := model.MetricInstance{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing data",
		})
	}

	// TODO: validation
	gron := gronx.New()
	if !gron.IsValid(payload.CronSchedule) {
		return c.Status(400).JSON(&fiber.Map{
			"message": "error validation cron_schedule",
		})
	}

	err := cc.DB.Transaction(func(tx *gorm.DB) error {
		queryResults := tx.Create(&payload).Error
		if queryResults != nil {
			return queryResults
		}

		// Create schedule
		command := entities.ProfileCommand{
			Logger:    cc.Logger,
			Db:        cc.DB,
			EventName: constants.EVENT_CREATE_METRIC_INSTANCE,
			Payload:   payload,
		}

		cc.SW.DataChan <- command
		return nil
	})
	// If rollback happened
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"message": "error creating metric_instance",
			"context": err.Error(),
		})
	}

	return c.Status(201).JSON(&payload)
}
