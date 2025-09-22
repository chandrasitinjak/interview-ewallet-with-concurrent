package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "root:password@tcp(127.0.0.1:3306)/ewallet_with_concurrent?parseTime=true"
	return sql.Open("mysql", dsn)
}
