package dao

import (
	"log"
	"mini-tiktok/app/entity"
)

func VideoAdd(video *entity.Publication) error {
	err := mysqlDB.AutoMigrate(&entity.Publication{})
	if err != nil {
		return err
	}
	if err := mysqlDB.Create(video).Error; err != nil {
		log.Println(err)
	}
	return err
}

func VideoGetByUser(userID int64, videoList *[]entity.Publication) error {
	err := mysqlDB.AutoMigrate(&entity.Publication{})
	if err != nil {
		return err
	}
	err = mysqlDB.Where("owner_id=?", userID).Find(&videoList).Error
	if err != nil {
		log.Println(err)
	}
	return err
}
