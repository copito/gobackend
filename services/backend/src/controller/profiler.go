package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/gateway"
	"github.com/copito/data_quality/src/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func scanRows(list *sql.Rows) (rows []map[string]interface{}) {
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

// func CreateExampleTask(ctx context.Context, eventKey string, payload int) {
// 	fmt.Printf("TEST EXAMPLE (%v) - time: %v\n", payload, time.Now().Format("2006-01-02 03:04:05"))
// }

func CreateProfilerTask(ctx context.Context, eventKey string, payload model.MetricInstance) {
	logger := ctx.Value("logger").(*slog.Logger)

	// try-catch
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if !ok {
				logger.Error("CreateProfilerTask failed to return a proper error", slog.Any("err", r))
			}
			logger.Error("CreateProfilerTask failed - recovery call", slog.String("err", err.Error()))
		}
	}()

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
		// snowflake://{user}:{password}@{account_identifier}/
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

	// Parse Host information based on connection string
	hostname, err := payload.DatabaseOnboarding.GetHostName()
	if err != nil {
		logger.Error("error parsing hostname for connection string", slog.String("err", err.Error()))
		return
	}

	// TODO: Render query (calculated_template)
	// TODO: Template
	templatedCalculation := payload.Metric.TemplatedCalculation
	logger.Info(
		"query template engine",
		slog.String("template", templatedCalculation[:200]),
		slog.String("metric_level", payload.Metric.MetricLevel),
		slog.String("metric_name", payload.Metric.Name),
	)

	var query string
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
		return
	}

	defer rows.Close()
	profileResults := scanRows(rows)
	logger.Info("profile result", slog.Any("result", profileResults))
	// First column is always the output for the profile result (named value column)

	// Prepare message to send to Kafka
	if len(profileResults) == 0 {
		logger.Error("no result from profiling result")
		return
	}

	profileResult := profileResults[0]
	valueOutput, ok := profileResult["value"] // assumes result are in the value column
	if !ok {
		logger.Error("no result from profiling result", slog.Any("profile_results", profileResults))
		return
	}
	valueOutputString := fmt.Sprintf("%v", valueOutput)

	kafkaMessageOutput := entities.ProfileEvent{
		Metric:         payload.Metric.Name,
		Level:          payload.Metric.MetricLevel,
		IsCustom:       payload.Metric.IsCustom,
		IsStandard:     payload.Metric.IsStandard,
		IsTemplated:    payload.Metric.IsTemplated,
		DatabaseType:   payload.DatabaseOnboarding.DBType,
		DatabaseHost:   hostname,
		DatabaseName:   *payload.DatabaseName,
		DatabaseSchema: *payload.SchemaName,
		DatabaseTable:  payload.TableName,
		Payload:        valueOutput,
		PayloadMarshal: valueOutputString,
		Tenancy:        constants.TENANCY_DEV,
	}
	value, err := json.Marshal(kafkaMessageOutput)
	if err != nil {
		panic(err)
	}

	// Send data to kafka
	topic := viper.GetString("kafka.profile_metric_topic")
	kgateway := gateway.NewKafkaGateway(ctx)
	kgateway.PublishResultToKafka(ctx, topic, eventKey, string(value))
}

func UpdateJobsBasedOnDatabase(ctx context.Context, db *gorm.DB, sw *entities.ScheduleWorker) {
	logger := ctx.Value("logger").(*slog.Logger)

	logger.Info("Starting to updated jobs based on Database")

	var wg sync.WaitGroup
	wg.Add(3) // Add, Update, Delete (3)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// TODO: Look through metric_instances on the database
		var addMetricInstances []model.MetricInstance
		model.GetAllNonRegisteredMetricInstances(db, &addMetricInstances)

		// TODO: Add any missing
		for _, addRequired := range addMetricInstances {
			command := entities.ProfileCommand{
				Logger:    logger,
				Db:        db,
				EventName: constants.EVENT_CREATE_METRIC_INSTANCE,
				Payload:   addRequired,
			}

			sw.DataChan <- command
		}
	}(&wg)

	// TODO: Update changed
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
	}(&wg)

	// TODO: Delete not used
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
	}(&wg)

	wg.Wait()
}
