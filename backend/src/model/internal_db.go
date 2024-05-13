package model

import (
	"github.com/copito/data_quality/src/constants"
)

// TABLE: database_onboarding
type DatabaseOnboarding struct {
	UUIDBaseModel // adds id, created_at, deleted_at, updated_at

	Name             string                 `gorm:"column:name;not null;unique" json:"name"`
	DBType           constants.DatabaseType `gorm:"column:db_type;not null" json:"db_type"` // snowflake,mysql,postgres
	ConnectionString string                 `gorm:"column:connection_string;not null" json:"connection_string"`

	CreatedBy string `gorm:"column:created_by;not null" json:"created_by"`
}
