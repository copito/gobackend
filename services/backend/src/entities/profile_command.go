package entities

import (
	"log/slog"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/model"
	"gorm.io/gorm"
)

type ProfileCommand struct {
	Logger *slog.Logger
	Db     *gorm.DB

	EventName constants.EventName
	Payload   model.MetricInstance
}

type ProfileEvent struct {
	Metric      string `json:"metric_name"`
	Level       string `json:"metric_level"`
	IsCustom    bool   `json:"is_custom"`
	IsStandard  bool   `json:"is_standard"`
	IsTemplated bool   `json:"is_templated"`

	DatabaseType   constants.DatabasePlatform `json:"database_type"`
	DatabaseHost   string                     `json:"database_host"`
	DatabaseName   string                     `json:"database_name"`
	DatabaseSchema string                     `json:"database_schema"`
	DatabaseTable  string                     `json:"database_table"`

	// Send message to topic
	Payload        interface{} `json:"payload"`
	PayloadMarshal string      `json:"payload_marshal"`

	// Environment
	Tenancy constants.Tenancy `json:"tenancy"`
}
