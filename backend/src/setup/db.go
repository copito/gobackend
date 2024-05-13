package setup

import (
	"log/slog"
	"strings"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase(logger *slog.Logger) *gorm.DB {
	logger.Info("Loading database configurations...")

	// Capture database configuration
	conn := viper.GetString("database.connection_string")
	if strings.Trim(conn, " ") == "" {
		panic("unable to read connection string")
	}

	// Connect to database
	var db *gorm.DB
	var errConnection error
	dbType := viper.GetString("database.type")
	switch dbType {
	case string(constants.SQLITE):
		// test.db
		db, errConnection = gorm.Open(sqlite.Open(conn), &gorm.Config{})
	case string(constants.POSTGRES):
		// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, errConnection = gorm.Open(postgres.Open(conn), &gorm.Config{})
	default:
		db, errConnection = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	}

	if errConnection != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&model.DatabaseOnboarding{},
		&model.Metric{},
		&model.MetricInstance{},
	)

	// Load initial data (metrics)
	LoadInitialDatabaseOnboarding(logger, db)
	LoadInitialMetrics(logger, db)
	LoadInitialMetricsInstance(logger, db)

	return db
}
