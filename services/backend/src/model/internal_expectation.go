package model

// Metric is a reference to a calculation using templated query that are saved to db
// TABLE: metric
type Expectation struct {
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
