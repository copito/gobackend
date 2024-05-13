package model

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Metric is a reference to a calculation using templated query that are saved to db
// TABLE: metric
type Metric struct {
	UUIDBaseModel        // adds id, created_at, deleted_at, updated_at
	Name          string `gorm:"column:name;not null" json:"name"`
	Description   string `gorm:"column:description" json:"description"`

	// Metadata
	MetricLevel string `gorm:"column:metric_level;not null" json:"metric_level"` // table, column
	IsStandard  bool   `gorm:"column:is_standard;not null" json:"is_standard"`   // is_standard: bool
	IsCustom    bool   `gorm:"column:is_custom;not null" json:"is_custom"`
	IsTemplated bool   `gorm:"column:is_templated;not null" json:"is_templated"`

	// Templated question to be run
	TemplatedCalculation string `gorm:"column:templated_calculation" json:"templated_calculation"` // query with at least table (to run against) and db_type (reference to onboarding type)
	AllowedDatabases     string `gorm:"column:allowed_databases" json:"allowed_databases"`         // []string (at least 1) => [mysql,postgres,snowflake] -> JSON
	Tags                 string `gorm:"column:tags" json:"tags"`                                   // map[string]string -> JSON

	// Creator information
	CreatedBy string `gorm:"column:created_by;not null" json:"created_by"`
}

func (m Metric) GetIdentifier(db *gorm.DB) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(m.MetricLevel), strings.ToLower(m.Name))
}

// MetricInstance is a combination of Database Table, Metric, Cron that will have
// a measurement in a database like InfluxDB
// TABLE: metric_instance
type MetricInstance struct {
	UUIDBaseModel // adds id, created_at, deleted_at, updated_at

	// Database
	// DatabaseOnboardingID uuid.UUID          `gorm:"column:database_onboarding_id"`
	DatabaseOnboardingID uint               `gorm:"column:database_onboarding_id" json:"database_onboarding_id"`
	DatabaseOnboarding   DatabaseOnboarding `gorm:"foreignKey:DatabaseOnboardingID"`

	SchemaName *string `gorm:"column:schema_name" json:"schema_name"`
	TableName  string  `gorm:"column:table_name;not null" json:"table_name"`
	Columns    *string `gorm:"column:columns" json:"columns"` // Serialized string representing columns (e.g., JSON)

	// Params to be passed to templated query
	// Params *string `gorm:"column:params" json:"params"` // JSON -> {'param1': "value1", etc} map[string]string

	// Metric
	// MetricID uuid.UUID `gorm:"column:metric_id"`
	MetricID uint   `gorm:"column:metric_id" json:"metric_id"`
	Metric   Metric `gorm:"foreignKey:MetricID"`

	// BaseSchedule
	BaseSchedule
}

func (m MetricInstance) GetIdentifier(db *gorm.DB) string {
	return fmt.Sprintf("%s_%s__%s", m.DatabaseOnboarding.DBType, m.TableName, m.Metric.GetIdentifier(db))
}
