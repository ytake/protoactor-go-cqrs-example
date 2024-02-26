package config

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	MySQLMaxOpenConns    = 30
	MySQLMaxIdleConns    = 40
	MySQLConnMaxLifetime = time.Second * 60
)

// MysqlConfig returns a mysql.Config object with the pre-configured settings
// for connecting to a MySQL database.
func MysqlConfig() *mysql.Config {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return &mysql.Config{
		DBName:    "sample",
		User:      "user",
		Passwd:    "passw@rd",
		Addr:      "localhost:3306",
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_bin",
		Loc:       jst,
	}
}
