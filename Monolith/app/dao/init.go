package dao

import (
	"mini-tiktok/common/db"

	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func init() {
	mysqlDB = db.NewMySQLConnInstance().DB
}
