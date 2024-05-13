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
