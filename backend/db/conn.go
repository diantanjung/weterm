package db

import (
	"database/sql"
	"fmt"

	"github.com/diantanjung/weterm/backend/util"
	_ "github.com/lib/pq"
)

func Open(config util.Config) (conn *sql.DB, err error) {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	conn, err = sql.Open(config.DBDriver, dbSource)
	return
}
