package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/model"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	createdByDefault = "data-platform-core"
)

type PreLoadedMetric struct {
	File             string
	Description      string
	IsTemplated      bool
	AllowedDatabases []string
	Tags             map[string]string
}

// breakDownFile breaks down the filename using double underscore into the metric name and the metric_level
func breakDownFile(fileName string) (name string, level string, err error) {
	fileName = strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	breakdown := strings.Split(fileName, "__")
	if len(breakdown) != 2 {
		err = errors.New("path does not follow the correct structure")
		return name, level, err
	}

	level = breakdown[0]
	name = breakdown[1]
	return name, level, nil
}

func getMetricDirectory(logger *slog.Logger) string {
	assetPath := "assets/metrics"
	sourcePath, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		return ""
	}

	path := path.Join(sourcePath, assetPath)
	logger.Debug("metrics paths built", slog.String("path", path))
	return path
}

func LoadInitialDatabaseOnboarding(logger *slog.Logger, db *gorm.DB) {
	// Check if dataset is empty (if not bypass)
	var count int64
	db.Debug().Model(&model.DatabaseOnboarding{}).Count(&count)

	if count > 0 {
		// Data already loaded
		logger.Info(
			"Data Already loaded - bypassing initial data",
			slog.Int64("data_onboarding_count", count),
		)
		return
	}

	dbType := viper.GetString("database.type")
	conn := viper.GetString("database.connection_string")
	result := db.Debug().Create(&model.DatabaseOnboarding{
		Name:             fmt.Sprintf("own_db_%v", uuid.NewString()),
		DBType:           constants.DatabaseType(dbType),
		ConnectionString: conn,
		CreatedBy:        createdByDefault,
	})

	logger.Debug(fmt.Sprintf("Result from pre-load: %v", result))
	logger.Info("Creating own database")
}

func LoadInitialMetrics(logger *slog.Logger, db *gorm.DB) {
	// Check if dataset is empty (if not bypass)
	var count int64
	db.Debug().Model(&model.Metric{}).Count(&count)

	if count > 0 {
		// Data already loaded
		logger.Info(
			"Data Already loaded - bypassing initial data",
			slog.Int64("metric_count", count),
		)
		return
	}

	defaultAllowedDatabases := []string{"postgres", "sqlite"}
	defaultTag := map[string]string{"pre_built": "True"}
	preBuildMetrics := []PreLoadedMetric{
		{
			File:             "table__row_count.sql",
			Description:      "Counts the number of rows in the table",
			AllowedDatabases: defaultAllowedDatabases,
			Tags:             defaultTag,
		},
		// {
		// 	File:             "table__column_count.sql",
		// 	Description:      "Counts the number of column in the table",
		// 	AllowedDatabases: defaultAllowedDatabases,
		// 	Tags:             defaultTag,
		// },
		// {
		// 	File:             "table__columns.sql",
		// 	Description:      "Evaluates to a list of columns in the table",
		// 	AllowedDatabases: defaultAllowedDatabases,
		// 	Tags:             defaultTag,
		// },
		// {
		// 	File:             "column__distinct_count.sql",
		// 	Description:      "Counts the distinct items in a column",
		// 	AllowedDatabases: defaultAllowedDatabases,
		// 	Tags:             defaultTag,
		// },
		// {
		// 	File:             "column__max.sql",
		// 	Description:      "Evaluates the maximum value for a column",
		// 	AllowedDatabases: defaultAllowedDatabases,
		// 	Tags:             defaultTag,
		// },
		{
			File:             "column__value_match_regex.sql",
			Description:      "Evaluate a column matches a certain regex",
			AllowedDatabases: defaultAllowedDatabases,
			Tags:             defaultTag,
			IsTemplated:      true,
		},
	}

	basePath := getMetricDirectory(logger)
	for _, metric := range preBuildMetrics {

		// Breakdown file name table__row_count.sql into table/row_count
		name, metricLevel, err := breakDownFile(metric.File)
		if err != nil {
			logger.Error("Unable to load metric", slog.String("file", metric.File))
			panic("error loading default metrics file")
		}

		// Check if valid metricLevel
		if !slices.Contains([]string{"table", "column"}, metricLevel) {
			logger.Error("Unable to load metric as metric level invalid", slog.String("file", metric.File))
			panic("error loading default metrics metric_level")
		}

		fullPath := path.Join(basePath, metric.File)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("Unable to load metric", slog.String("file_path", fullPath), slog.String("issue", "templated_calculation"))
			panic("error loading default metrics file data")
		}

		allowedDatabasesBytes, err := json.Marshal(metric.AllowedDatabases)
		if err != nil {
			logger.Error("Unable to load metric allowed_databases", slog.String("file_path", fullPath), slog.String("issue", "allowed_databases"))
		}

		tagsBytes, err := json.Marshal(metric.Tags)
		if err != nil {
			logger.Error("Unable to load metric", slog.String("file_path", fullPath), slog.String("issue", "tags"))
			panic("error loading default metrics tags")
		}

		result := db.Debug().Create(&model.Metric{
			Name:                 name,
			MetricLevel:          metricLevel,
			Description:          metric.Description,
			IsStandard:           true,
			IsCustom:             false,
			IsTemplated:          metric.IsTemplated,
			TemplatedCalculation: string(data),
			AllowedDatabases:     string(allowedDatabasesBytes),
			Tags:                 string(tagsBytes),
			CreatedBy:            createdByDefault,
		})

		logger.Debug(fmt.Sprintf("Result from pre-load: %v", result))
		logger.Info(
			"Creating standard metric row in database",
			slog.String("name", name),
			slog.String("metric_level", metricLevel),
			// slog.Int64("rows_affected", result.RowsAffected),
			// slog.String("err", result.Error.Error()),
		)
	}
}

func LoadInitialMetricsInstance(logger *slog.Logger, db *gorm.DB) {
	// Check if dataset is empty (if not bypass)
	var count int64
	db.Debug().Model(&model.MetricInstance{}).Count(&count)

	if count > 0 {
		// Data already loaded
		logger.Info(
			"Data Already loaded - bypassing initial data",
			slog.Int64("metric_instance_count", count),
		)
		return
	}

	defaultDB := viper.Get("database.type")
	if defaultDB == "" {
		panic("database type not provided in configs: database.type")
	}

	var dbResult model.DatabaseOnboarding
	db.Debug().Model(&model.DatabaseOnboarding{}).Where("db_type = ?", defaultDB).Find(&dbResult)

	var metricRowCountResult model.Metric
	db.Debug().Model(&model.Metric{}).Where("name = ?", "row_count").Find(&metricRowCountResult)

	var metricValueMatchRegexResult model.Metric
	db.Debug().Model(&model.Metric{}).Where("name = ?", "value_match_regex").Find(&metricValueMatchRegexResult)

	dbName := "postgres"
	dbSchama := "public"
	result := db.Debug().Create(&model.MetricInstance{
		DatabaseOnboardingID: dbResult.ID,
		DatabaseName:         &dbName,
		SchemaName:           &dbSchama,
		TableName:            "metrics",
		Columns:              nil,
		MetricID:             metricRowCountResult.ID,
		BaseSchedule: model.BaseSchedule{
			CronSchedule:    "*/1 * * * *", // Every 1 minute
			ScheduleJobID:   nil,
			ScheduleGateway: constants.GOCRON,
		},
	})

	logger.Debug(fmt.Sprintf("Result from pre-load: %v", result))
	logger.Info("Creating standard metric instance row in database")

	result2 := db.Debug().Create(&model.MetricInstance{
		DatabaseOnboardingID: dbResult.ID,
		DatabaseName:         &dbName,
		SchemaName:           &dbSchama,
		TableName:            "metrics",
		Columns:              nil,
		MetricID:             metricValueMatchRegexResult.ID,
		BaseSchedule: model.BaseSchedule{
			CronSchedule:    "5 4 * * *", // At 04:05 every day
			ScheduleJobID:   nil,
			ScheduleGateway: constants.GOCRON,
		},
	})

	logger.Debug(fmt.Sprintf("Result from pre-load: %v", result2))
	logger.Info("Creating standard metric instance row in database")
}
