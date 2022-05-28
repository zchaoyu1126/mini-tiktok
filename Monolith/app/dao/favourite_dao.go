package dao

import (
	"log"
	"mini-tiktok/app/entity"
)

func FavouriteAdd(favourite *entity.Favourite) error {
	mysqlDB.AutoMigrate(&entity.Favourite{})
	err := mysqlDB.Create(favourite).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func FavouriteGetByVideoUser(favourite *entity.Favourite, uid int64, vid int64) error {
	mysqlDB.AutoMigrate(&entity.Favourite{})
	err := mysqlDB.Where("user_id=?", uid).Where("video_id=?", vid).Find(&favourite).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateFavourite(favourite *entity.Favourite) error {
	mysqlDB.AutoMigrate(&entity.Favourite{})
	err := mysqlDB.Where("user_id=?", favourite.UserID).Where("video_id=?", favourite.VideoID).Update("is_favourite", favourite.IsFavourite).Error
	if err != nil {
		log.Println(err)
	}
	return err
}
