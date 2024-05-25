package constants

type (
	DatabasePlatform string
)

const (
	MYSQL     DatabasePlatform = "mysql"
	POSTGRES  DatabasePlatform = "postgres"
	SQLITE    DatabasePlatform = "sqlite"
	SQLSERVER DatabasePlatform = "mssql"
	SNOWFLAKE DatabasePlatform = "snowflake"
	MARIADB   DatabasePlatform = "mariadb"
	DUCKDB    DatabasePlatform = "duckdb"
	BIGQUERY  DatabasePlatform = "bigquery"
)
