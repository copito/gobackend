package constants

type (
	Tenancy string
)

const (
	TENANCY_PROD Tenancy = "company/production"
	TENANCY_UAT  Tenancy = "company/uat"
	TENANCY_DEV  Tenancy = "company/dev"
)
