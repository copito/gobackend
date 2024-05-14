package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/copito/data_quality/src/controller"
	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/setup"
	"github.com/spf13/viper"
)

func main() {
	// Setup Logging
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelWarn)
	lvl.Set(slog.LevelDebug)
	logger := setup.SetupLogging(true, lvl.Level())

	// Setup all configuration for this application
	setup.SetupConfig(logger)
	db := setup.SetupDatabase(logger)

	if db == nil {
		logger.Error("closing application gracefully due to non establishing db connection")
		return
	}

	// Create a new context with the logger attached
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)

	// Start Cron Job Worker
	// TODO: change to a proper cronjob later
	doneChan := make(chan bool)
	dataChan := make(chan entities.ProfileCommand)
	scheduleWorker := entities.ScheduleWorker{
		DoneChan: doneChan,
		DataChan: dataChan,
	}
	// Create worker code
	go controller.CreateScheduleWorker(ctx, db, &scheduleWorker)

	// Update jobs
	go controller.UpdateJobsBasedOnDatabase(ctx, db, &scheduleWorker)

	// Setup Fiber (TODO: move to use uberfx - dependency injection)
	application := setup.NewApplication(ctx, logger, db, &scheduleWorker)
	application = application.SetupAPI()
	application = application.SetupMiddleware()
	application = application.SetupRoutes()

	// Force Stop at any point (with termination signals)
	// handled by fx too
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, os.Interrupt)

	go func(logger *slog.Logger, sig chan os.Signal, app *setup.Application) {
		<-sig
		logger.Warn("Forced exited the cli...")
		app.ScheduleWorker.DoneChan <- true
		logger.Warn("Gracefully shutting down cron server...")
		_ = app.Api.Shutdown()
		logger.Warn("Gracefully shutting down api server...")
		os.Exit(1)
	}(logger, sig, application)

	port := viper.GetString("backend.port")
	if port == "" {
		logger.Error("no port provided in configuration...")
		panic("no backend.port provided in configs")
	}

	err := application.Api.Listen(port)
	if err != nil {
		logger.Error("application failed...", slog.String("err", err.Error()))
		return
	}

	select {}
}
