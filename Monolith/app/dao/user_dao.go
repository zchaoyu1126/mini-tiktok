package dao

import (
	"log"
	"mini-tiktok/app/entity"
)

func UserAdd(user *entity.User) error {
	// migrate 仅支持创建表、增加表中没有的字段和索引。为了保护你的数据，它并不支持改变已有的字段类型或删除未被使用的字段
	mysqlDB.AutoMigrate(&entity.User{})
	err := mysqlDB.Create(user).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func UserGetByName(user *entity.User) error {
	mysqlDB.AutoMigrate(&entity.User{})
	err := mysqlDB.Where("username=?", user.UserName).First(user).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func UserGetByUID(user *entity.User) error {
	mysqlDB.AutoMigrate(&entity.User{})
	err := mysqlDB.Where("user_id=?", user.UserID).First(user).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func UserUpdate(user *entity.User) error {
	err := mysqlDB.Save(user).Error
	if err != nil {
		log.Println(err)
	}
	return err
}
