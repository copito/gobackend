package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/entities"
	"github.com/copito/data_quality/src/model"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func determineCron(cron string) gocron.JobDefinition {
	breakDownCron := strings.Split(cron, " ")
	if len(breakDownCron) == 5 {
		return gocron.CronJob(cron, false)
	} else if len(breakDownCron) == 6 {
		return gocron.CronJob(cron, true)
	} else {
		panic("incorrect cron structure")
	}
}

func CreateScheduleWorker(ctx context.Context, db *gorm.DB, sw *entities.ScheduleWorker) {
	logger := ctx.Value("logger").(*slog.Logger)

	doneChan := sw.DoneChan
	dataChan := sw.DataChan

	// create a scheduler
	s, err := gocron.NewScheduler()
	// s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		panic("Unable to start Profiler Worker")
	}

	// Start running scheduler
	s.Start()

readChannel:
	for {
		select {
		case <-doneChan:
			logger.Info("Closing down CRON...")
			s.Shutdown()
			break readChannel

		case event := <-dataChan:

			switch event.EventName {
			case constants.EVENT_CLEAN_ALL_METRIC_INSTANCE:
				logger.Warn("cleaning up all scheduled jobs")
				db.Transaction(func(tx *gorm.DB) error {
					// Update cron job id
					err = tx.Debug().Model(&model.MetricInstance{}).Update("schedule_job_id", nil).Error
					if err != nil {
						// Transaction failed and rolled back push to database
						// Must clean up scheduler
						// (ignore error as if it errors out it means job was not created)
						logger.Error("unable to clean up job ids in database")
						return err
					}

					s.RemoveByTags(string(constants.TAG_PROFILER))
					logger.Info("cleaned up all jobs with profiler as a tag")
					// add a job to the scheduler

					return nil
				})

			case constants.EVENT_CREATE_METRIC_INSTANCE:
				// This will be the information MetricInstance that was just created
				logger.Info("received create event", slog.Any("event", event))

				if event.Payload.ScheduleGateway != constants.GOCRON {
					logger.Warn("scheduling only implemented for gocron", "implemented", false)
					break
				}

				eventKey := fmt.Sprintf("%v.%v.%v", event.EventName, event.Payload.Metric.MetricLevel, event.Payload.Metric.Name)
				jobDefinition := determineCron(event.Payload.CronSchedule)

				// Error handled idempotently (self-clean up)
				err = db.Transaction(func(tx *gorm.DB) error {
					// add a job to the scheduler
					job, err := s.NewJob(
						jobDefinition, // -> gocron.CronJob(event.Payload.CronSchedule, false),
						gocron.NewTask(CreateProfilerTask, ctx, eventKey, event.Payload),
						gocron.WithTags(string(constants.TAG_PROFILER)), // add tags
						gocron.WithName(string(event.EventName)),        // provide name
						// gocron.JobOption(gocron.WithStartImmediately()), // start profiler just now
					)
					// Check if the new job can be created successfully
					if err != nil {
						// handle error
						logger.Error(
							"unable to create new job",
							slog.String("event_name", string(event.EventName)),
							slog.String("combined", eventKey),
							slog.String("err", err.Error()),
						)
						return err
					}
					// each job has a unique id
					logger.Info(
						"job has been created",
						slog.String("event_name", string(event.EventName)),
						slog.String("combined", eventKey),
						slog.String("id", job.ID().String()),
					)

					// Update cron job id
					err = tx.Debug().Model(&model.MetricInstance{}).Where("id = ?", event.Payload.ID).Update("schedule_job_id", job.ID().String()).Error
					if err != nil {
						// Transaction failed and rolled back push to database
						// Must clean up scheduler
						// (ignore error as if it errors out it means job was not created)
						s.RemoveJob(job.ID())
						return err
					}

					return nil
				})
				// If rollback occurs logging what was the cause
				if err != nil {
					logger.Error("rolledback metric instance creation", slog.String("err", err.Error()))
				}

			case constants.EVENT_DELETE_METRIC_INSTANCE:
				logger.Info("deleting metric_instance from cron", "implemented", false)
				jobID := event.Payload.ScheduleJobID

				if jobID == nil || *jobID == "" {
					logger.Error("schedule job id is empty")
					continue
				}

				jobUUID, err := uuid.Parse(*jobID)
				if err != nil {
					logger.Error("unable to parse job uuid for schedule", slog.Any("job_id", *jobID))
					continue
				}

				s.RemoveJob(jobUUID)

			case constants.EVENT_UPDATE_METRIC_INSTANCE:
				// jobDefinition := DetermineCron(event.Payload.CronSchedule)
				// s.Update(
				// 	"uuid",
				// 	jobDefinition, // -> gocron.CronJob(event.Payload.CronSchedule, false),
				// 	gocron.NewTask(CreateProfilerTask, "taskid", "hello", event.Payload),
				// 	gocron.WithTags("profiler"),              // add tags
				// 	gocron.WithName(string(event.EventName)), // provide name
				// 	// gocron.JobOption(gocron.WithStartImmediately()), // start profiler just now
				// )
				logger.Info("updating metric_instance from cron", "not_implemented", true)

			default:
			}
			// DEBUG: Add this one after debug
			// case <-time.After(5 * time.Minute):
			// 	// debugging purposes
			// 	fmt.Println("Timeout reached - closing down scheduler")
			// 	break readChannel
		}
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
		fmt.Println(err)
	}
}
