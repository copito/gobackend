package model

import (
	"time"

	"github.com/copito/data_quality/src/constants"
)

// Base model to be used in every model
type UUIDBaseModel struct {
	// ID        uuid.UUID  `gorm:"primary_key;unique;type:uuid;column:id;default:uuid_generate_v4()"`
	ID        uint       `gorm:"primaryKey;unique;column:id" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// Base model to be used for scheduler
// BaseSchedule structure base
type BaseSchedule struct {
	CronSchedule    string                     `gorm:"column:cron_schedule" json:"cron_schedule"`
	ScheduleJobID   *string                    `gorm:"column:schedule_job_id" json:"schedule_job_id"`
	ScheduleGateway constants.SchedulerGateway `gorm:"column:schedule_gateway" json:"schedule_gateway"`
}
