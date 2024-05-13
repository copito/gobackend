package constants

type EventName string

const (
	CREATE_METRIC_INSTANCE EventName = "metric_instance.create"
	UPDATE_METRIC_INSTANCE EventName = "metric_instance.update"
	DELETE_METRIC_INSTANCE EventName = "metric_instance.delete"
)
