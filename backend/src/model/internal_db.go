package model

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/copito/data_quality/src/constants"
)

// TABLE: database_onboarding
type DatabaseOnboarding struct {
	UUIDBaseModel // adds id, created_at, deleted_at, updated_at

	Name             string                 `gorm:"column:name;not null;unique" json:"name"`
	DBType           constants.DatabaseType `gorm:"column:db_type;not null" json:"db_type"` // snowflake,mysql,postgres
	ConnectionString string                 `gorm:"column:connection_string;not null" json:"connection_string"`

	CreatedBy string `gorm:"column:created_by;not null" json:"created_by"`
}

func (d DatabaseOnboarding) GetHostName() (string, error) {
	switch d.DBType {
	case constants.SQLITE:
		return d.ConnectionString, nil
	case
		constants.SQLSERVER,
		constants.POSTGRES,
		constants.MARIADB,
		constants.MYSQL:
		// "sqlserver://sa:12345678@localhost:1433?database=gorm"
		// Parse Host information based on connection string
		u, err := url.Parse(d.ConnectionString)
		if err != nil {
			fmt.Println(err)
			return "", errors.New("unable to parse connection string")
		}

		// Get the hostname
		hostname := u.Host
		return hostname, nil

	default:
		// Parse Host information based on connection string
		u, err := url.Parse(d.ConnectionString)
		if err != nil {
			fmt.Println(err)
			return "", errors.New("unable to parse connection string")
		}

		// Get the hostname
		hostname := u.Host
		return hostname, nil
	}
}
