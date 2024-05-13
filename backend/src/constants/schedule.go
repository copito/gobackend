package constants

type SchedulerGateway string

const (
	GOCRON              SchedulerGateway = "gocron"
	KUBERNETES_CRON_JOB SchedulerGateway = "kubernetes_cron_job"
	CADENCE             SchedulerGateway = "cadence"
	APACHE_AIRFLOW      SchedulerGateway = "apache_airflow"
)
