package dao

import (
	"log"
	"mini-tiktok/app/entity"
)

func FollowCreate(follow *entity.Follow) error {
	mysqlDB.AutoMigrate(&entity.Follow{})
	err := mysqlDB.Create(follow).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func FollowGetByIDs(follow *entity.Follow) error {
	mysqlDB.AutoMigrate(&entity.Follow{})
	err := mysqlDB.Where("from_user_id = ? AND to_user_id = ?", follow.FromUserID, follow.ToUserID).First(follow).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func FollowUpdate(follow *entity.Follow) error {
	mysqlDB.AutoMigrate(&entity.Follow{})
	err := mysqlDB.Save(follow).Error
	if err != nil {
		log.Println(err)
	}
	return err
}
