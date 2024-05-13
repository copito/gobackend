package controller

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Controller struct {
	Context context.Context
	Logger  *slog.Logger
	DB      *gorm.DB
	Api     *fiber.App
}

func NewController(ctx context.Context, logger *slog.Logger, app *fiber.App, db *gorm.DB) *Controller {
	return &Controller{
		Context: ctx,
		Logger:  logger,
		Api:     app,
		DB:      db,
	}
}
