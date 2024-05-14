package entities

import "github.com/copito/data_quality/src/constants"

type SchedulingConfig struct {
	CronScheduleInSeconds string
	ScheduleGateway       constants.SchedulerGateway

	// TODO: when to start/etc/gocron options
}

type ScheduleWorker struct {
	DoneChan chan bool
	DataChan chan ProfileCommand
}
