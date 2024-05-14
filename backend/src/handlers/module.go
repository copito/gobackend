package handlers

import (
	"context"
	"log/slog"

	"github.com/copito/data_quality/src/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2/middleware/monitor"
)

type Handlers struct {
	Context context.Context
	Logger  *slog.Logger
	DB      *gorm.DB
	Api     *fiber.App
	SW      *entities.ScheduleWorker
}

func NewHandlers(ctx context.Context, logger *slog.Logger, app *fiber.App, db *gorm.DB, sw *entities.ScheduleWorker) *Handlers {
	return &Handlers{
		Context: ctx,
		Logger:  logger,
		Api:     app,
		DB:      db,
		SW:      sw,
	}
}

func (h *Handlers) LoadEndpoints() *fiber.App {
	h.Api.Get("/", func(c *fiber.Ctx) error {
		// Render a template named 'index.html' with content
		return c.Render("index", fiber.Map{
			"Title":       "Hello, World!",
			"Description": "This is a template.",
		})
	})

	h.Api.Get("/health/", h.HealthCheck)
	h.Api.Get("/stack/", h.Stack)
	h.Api.Get("/monitor/", monitor.New(monitor.Config{Title: "Metrics Page"}))

	api := h.Api.Group("/api/")
	apiV1 := api.Group("/v1/")

	apiV1.Get("/db_onboard/", h.GetDatabaseOnboard)
	apiV1.Post("/db_onboard/", h.CreateDatabaseOnboard)
	apiV1.Get("/db_onboard/:id", h.GetDatabaseOnboardByID)

	apiV1.Get("/metrics/", h.GetMetrics)
	apiV1.Post("/metrics/", h.CreateMetric)
	apiV1.Get("/metrics/:id", h.GetMetricByID)

	apiV1.Get("/metric_instances/", h.GetMetricInstances)
	apiV1.Post("/metric_instances/", h.CreateMetricInstanceByID)
	apiV1.Get("/metric_instances/:id/", h.GetMetricInstanceByID)

	return h.Api
}
