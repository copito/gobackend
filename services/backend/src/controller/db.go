package controller

import (
	"github.com/copito/data_quality/src/constants"
	"github.com/copito/data_quality/src/model"
)

// TODO: NotImplemented yet need implement
func CheckDatabaseConnectivity(db model.DatabaseOnboarding) (bool, error) {
	isConnected := false

	switch db.DBType {
	case constants.SQLITE:
	case constants.MYSQL:
	case constants.POSTGRES:
	case constants.SQLSERVER:
	case constants.MARIADB:
	case constants.DUCKDB:
	case constants.SNOWFLAKE:
	case constants.BIGQUERY:
	default:

	}

	return isConnected, nil
}
