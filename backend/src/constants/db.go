package constants

type (
	DatabaseType string
)

const (
	MYSQL     DatabaseType = "mysql"
	POSTGRES  DatabaseType = "postgres"
	SQLITE    DatabaseType = "sqlite"
	SQLSERVER DatabaseType = "mssql"
	SNOWFLAKE DatabaseType = "snowflake"
	MARIADB   DatabaseType = "mariadb"
	DUCKDB    DatabaseType = "duckdb"
	BIGQUERY  DatabaseType = "bigquery"
)
