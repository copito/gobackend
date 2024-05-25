package constants

type EventName string

const (
	EVENT_CLEAN_ALL_METRIC_INSTANCE EventName = "metric_instance.force_clean"
	EVENT_CREATE_METRIC_INSTANCE    EventName = "metric_instance.create"
	EVENT_UPDATE_METRIC_INSTANCE    EventName = "metric_instance.update"
	EVENT_DELETE_METRIC_INSTANCE    EventName = "metric_instance.delete"
)

type Tag string

const (
	TAG_PROFILER    Tag = "profiler"
	TAG_EXPECTATION Tag = "expectation"
)
