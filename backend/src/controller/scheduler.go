package controller

type Workflower interface {
	Register(cron string) (string, error)
	Run() error
}

func ScheduleWorkflow() {
}
