package setup

import (
	"context"
	"log/slog"
	"time"

	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/gofiber/fiber/v2/middleware/healthcheck"
)

type Application struct {
	Context        context.Context
	Logger         *slog.Logger
	DB             *gorm.DB
	Api            *fiber.App
	ScheduleWorker *entities.ScheduleWorker
}

func NewApplication(ctx context.Context, logger *slog.Logger, db *gorm.DB, scheduleWorker *entities.ScheduleWorker) *Application {
	return &Application{
		Context:        ctx,
		Logger:         logger,
		Api:            nil,
		DB:             db,
		ScheduleWorker: scheduleWorker,
	}
}

func (a *Application) SetupAPI() *Application {
	a.Logger.Info("Setting Up API...")

	app := fiber.New()
	app.Static("/frontend/", "./../../frontend", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	a.Api = app
	return a
}

func (a *Application) SetupRoutes() *Application {
	a.Logger.Info("Setting Up API Routes...")

	// Build controller
	cc := handlers.NewHandlers(a.Context, a.Logger, a.Api, a.DB, a.ScheduleWorker)
	cc.LoadEndpoints()

	return a
}

func (a *Application) SetupMiddleware() *Application {
	a.Logger.Info("Setting Up API Middleware...")

	// Cors
	a.Api.Use(cors.New())

	// 20 requests per 10 seconds max
	a.Api.Use(limiter.New(limiter.Config{
		Expiration:         10 * time.Second,
		Max:                20,
		SkipFailedRequests: true,
	}))

	// Protect against panics (helping with recovery)
	a.Api.Use(recover.New())

	// Health check
	// app.Use(healthcheck.New(healthcheck.Config{
	// 	LivenessProbe: func(c *fiber.Ctx) bool {
	// 		return true
	// 	},
	// 	LivenessEndpoint: "/live",
	// 	ReadinessProbe: func(c *fiber.Ctx) bool {
	// 		return serviceA.Ready() && serviceB.Ready() && ...
	// 	},
	// 	ReadinessEndpoint: "/ready",
	// }))

	return a
}
