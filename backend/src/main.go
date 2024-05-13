package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/controller"
	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/model"
	"github.com/copito/data_quality/src/setup"
	"github.com/spf13/viper"
)

func main() {
	// Force Stop at any point (with termination signals)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	go func(sig chan os.Signal) {
		<-sig
		fmt.Println("Forced exited the cli...")
		os.Exit(1)
	}(sig)

	// Setup Logging
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelWarn)
	lvl.Set(slog.LevelDebug)
	logger := setup.SetupLogging(true, lvl.Level())

	// Setup all configuration for this application
	setup.SetupConfig(logger)
	db := setup.SetupDatabase(logger)

	// Create a new context with the logger attached
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)

	// Start Cron Job Worker
	// TODO: change to a proper cronjob later
	doneChan := make(chan bool)
	dataChan := make(chan entities.ProfileCommand)
	scheduleWorker := entities.ScheduleWorker{
		doneChan: doneChan,
		dataChan: dataChan,
	}
	go controller.CreateProfilerWorker(ctx, db, doneChan, dataChan)

	// DEBUG: test dynamic creating
	// DEBUG: test dynamic creating
	// DEBUG: test dynamic creating

	var metricInstanceResult model.MetricInstance
	db.Debug().Preload("DatabaseOnboarding").Preload("Metric").Model(&model.MetricInstance{}).Find(&metricInstanceResult, 1)

	command := entities.ProfileCommand{
		Logger: logger, Db: db, EventName: constants.CREATE_METRIC_INSTANCE,
		Payload: metricInstanceResult,
	}
	// command2 := entities.ProfileCommand{
	// 	Logger: logger, Db: db, EventName: constants.CREATE_METRIC_INSTANCE,
	// 	Payload: metricInstanceResult,
	// }

	time.Sleep(5 * time.Second)
	dataChan <- command
	time.Sleep(20 * time.Minute)
	// dataChan <- command2

	// CONTINUE
	// Setup Fiber (TODO: move to use uberfx - dependency injection)
	application := setup.NewApplication(ctx, logger, db, scheduleWorker)
	application = application.SetupAPI()
	application = application.SetupMiddleware()
	application = application.SetupRoutes()

	port := viper.GetString("backend.port")
	if port == "" {
		logger.Error("no port provided in configuration...")
		panic("no backend.port provided in configs")
	}

	err := application.Api.Listen(port)
	if err != nil {
		logger.Error("application failed...", slog.String("err", err.Error()))
	}

	select {}
}
