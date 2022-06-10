package dao

import (
	"mini-tiktok/app/entity"
	"mini-tiktok/common/db"

	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func init() {
	mysqlDB = db.NewMySQLConnInstance().DB
	mysqlDB.AutoMigrate(&entity.User{})
	mysqlDB.AutoMigrate(&entity.Publication{})
	mysqlDB.AutoMigrate(&entity.Favourite{})
	mysqlDB.AutoMigrate(&entity.Follow{})
	mysqlDB.AutoMigrate(&entity.Comment{})
}
