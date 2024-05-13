package controller

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/gateway"
	"github.com/copito/data_quality/src/model"
	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateProfilerWorker(ctx context.Context, db *gorm.DB, doneChan chan bool, dataChan chan entities.ProfileCommand) {
	logger := ctx.Value("logger").(*slog.Logger)

	// create a scheduler
	s, err := gocron.NewScheduler()
	// s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		panic("Unable to start Profiler Worker")
	}

	// Start running scheduler
	s.Start()

readChannel:
	for {
		select {
		case <-doneChan:
			logger.Info("Closing down CRON...")
			break readChannel

		case event := <-dataChan:

			switch event.EventName {
			case constants.CREATE_METRIC_INSTANCE:
				// This will be the information MetricInstance that was just created
				// TODO: profiler task
				logger.Info("received event", slog.Any("event", event))

				if event.Payload.ScheduleGateway != constants.GOCRON {
					logger.Warn("scheduling only implemented for gocron", "implemented", false)
					break
				}

				eventKey := fmt.Sprintf("%v.%v.%v", event.EventName, event.Payload.Metric.MetricLevel, event.Payload.Metric.Name)

				var job gocron.Job
				err = db.Transaction(func(tx *gorm.DB) error {
					// add a job to the scheduler
					job, err := s.NewJob(
						gocron.CronJob(event.Payload.CronSchedule, true), // schedule (from payload)
						gocron.NewTask(CreateProfilerTask, ctx, eventKey, event.Payload),
						gocron.WithTags("profiler"),              // add tags
						gocron.WithName(string(event.EventName)), // provide name
						// gocron.JobOption(gocron.WithStartImmediately()), // start profiler just now
					)
					// Check if the new job can be created successfully
					if err != nil {
						// handle error
						logger.Error(
							"unable to create new job",
							slog.String("event_name", string(event.EventName)),
							slog.String("combined", eventKey),
							slog.String("err", err.Error()),
						)
						return err
					}
					// each job has a unique id
					logger.Info(
						"job has been created",
						slog.String("event_name", string(event.EventName)),
						slog.String("combined", eventKey),
						slog.String("id", job.ID().String()),
					)

					// Update cron job id
					err = tx.Debug().Model(&model.MetricInstance{}).Where("id = ?", event.Payload.ID).Update("schedule_job_id", job.ID().String()).Error
					if err != nil {
						return err
					}

					return nil
				})
				// Transaction failed and rolled back push to database
				// Must clean up scheduler
				if err != nil {
					s.RemoveJob(job.ID())
				}

			case constants.DELETE_METRIC_INSTANCE:
				// s.RemoveJob("uuid")
				logger.Info("deleting metric_instance from cron", "implemented", false)

			case constants.UPDATE_METRIC_INSTANCE:
				// s.Update(
				// 	"uuid",
				// 	gocron.CronJob(event.Scheduling.CronScheduleInSeconds, true), // schedule (from payload),
				// 	gocron.NewTask(CreateProfilerTask, "taskid", "hello", event.Payload),
				// 	gocron.WithTags("profiler"),              // add tags
				// 	gocron.WithName(string(event.EventName)), // provide name
				// 	// gocron.JobOption(gocron.WithStartImmediately()), // start profiler just now
				// )
				logger.Info("updating metric_instance from cron", "implemented", false)

			default:
			}
		// DEBUG: Remove this one after debug
		case <-time.After(5 * time.Minute):
			// debugging purposes
			fmt.Println("Timeout reached - closing down scheduler")
			break readChannel
		}
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
		fmt.Println(err)
	}
}

func gormDatabaseToSQLDatabase(logger *slog.Logger, db *gorm.DB, errConnection error) *sql.DB {
	// Check connection
	if errConnection != nil {
		logger.Error("error opening connection", slog.String("error", errConnection.Error()))
		panic("unable to open a connection")
	}

	dd, err := db.DB()
	if err != nil {
		logger.Error("error converting GORM connection to plain sql connection", slog.String("error", err.Error()))
		panic("error converting GORM connection to plain sql connection")
	}

	return dd
}

func ScanRows(list *sql.Rows) (rows []map[string]interface{}) {
	fields, _ := list.Columns()
	for list.Next() {
		scans := make([]interface{}, len(fields))
		row := make(map[string]interface{})

		for i := range scans {
			scans[i] = &scans[i]
		}
		list.Scan(scans...)
		for i, v := range scans {
			value := ""
			if v != nil {
				value = fmt.Sprintf("%s", v)
			}
			row[fields[i]] = value
		}
		rows = append(rows, row)
	}
	return rows
}

func CreateProfilerTaskExample(ctx context.Context, eventKey string, payload int) {
	fmt.Printf("TEST EXAMPLE (%v) - time: %v\n", payload, time.Now().Format("2006-01-02 03:04:05"))
}

func CreateProfilerTask(ctx context.Context, eventKey string, payload model.MetricInstance) {
	logger := ctx.Value("logger").(*slog.Logger)

	// Preloaded Metric & DatabaseOnboarding
	// parse metric_instance data
	logger.Info(
		"running profile task",
		slog.String("event_key", eventKey),
		slog.String("tenancy", "production"),
	)

	// open connection to database
	conn := payload.DatabaseOnboarding.ConnectionString
	var dd *gorm.DB
	var db *sql.DB
	var errConnection error
	dbType := viper.GetString("database.type")
	switch dbType {
	case string(constants.SQLITE):
		// import "gorm.io/driver/sqlite" // Sqlite driver based on CGO
		// import "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
		// test.db
		dd, errConnection = gorm.Open(sqlite.Open(conn), &gorm.Config{})
		db = gormDatabaseToSQLDatabase(logger, dd, errConnection)
		defer db.Close()
	case string(constants.POSTGRES):
		// import "gorm.io/driver/postgres"
		// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		dd, errConnection = gorm.Open(postgres.Open(conn), &gorm.Config{})
		db = gormDatabaseToSQLDatabase(logger, dd, errConnection)
		defer db.Close()
	case string(constants.SQLSERVER):
		// import "gorm.io/driver/sqlserver"
		// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		dd, errConnection = gorm.Open(postgres.Open(conn), &gorm.Config{})
		db = gormDatabaseToSQLDatabase(logger, dd, errConnection)
		defer db.Close()
	case string(constants.MYSQL):
		// import "gorm.io/driver/mysql"
		// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		dd, errConnection = gorm.Open(postgres.Open(conn), &gorm.Config{})
		db = gormDatabaseToSQLDatabase(logger, dd, errConnection)
		defer db.Close()
	case string(constants.SNOWFLAKE):
		// import _ "github.com/snowflakedb/gosnowflake"
		// "user:password@my_organization-my_account/mydb"
		db, errConnection = sql.Open("snowflake", conn)
		defer db.Close()
		if errConnection != nil {
			logger.Error("unable to connect to snowflake database", slog.String("err", errConnection.Error()))
			panic("error connecting to snowflake")
		}
	default:
		dd, errConnection = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		db = gormDatabaseToSQLDatabase(logger, dd, errConnection)
		defer db.Close()
	}

	// TODO: Render query (calculated_template)
	var query string

	// TODO: Template
	templatedCalculation := payload.Metric.TemplatedCalculation
	logger.Info(
		"query template engine",
		slog.String("template", templatedCalculation[:200]),
		slog.String("metric_level", payload.Metric.MetricLevel),
		slog.String("metric_name", payload.Metric.Name),
	)

	if payload.Metric.IsTemplated {
		// TODO: pass more than database metadata, pass value too
		// Use payload.Params - use to template
		// passed params: database_type, table_name, param1=value1, param2=value2,...
		query = "SELECT COUNT(*) as value FROM metrics;"
	} else {
		// passed params: database_type, table_name
		query = "SELECT COUNT(*) as value FROM metrics;"
	}

	// TODO: Run query (get result)
	rows, err := db.Query(query)
	if err != nil {
		logger.Error("error running profiler query", slog.String("truncated_query", query[:200]))
		panic("error running profiler query")
	}

	defer rows.Close()
	profileResult := ScanRows(rows)
	logger.Info("profile result", slog.Any("result", profileResult))
	// First column is always the output for the profile result (named value column)

	// // TODO: send data to kafka
	// // Send data to kafka
	gateway.PublishResultToKafka(ctx, profileResult)
}
