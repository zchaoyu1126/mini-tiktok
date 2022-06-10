package dao

import (
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
)

func FavouriteAdd(favourite *entity.Favourite) error {
	err := mysqlDB.Create(favourite).Error
	if err != nil {
		zap.L().Error("mysql:favourites add new entry failed")
		return xerr.ErrDatabase
	}
	return nil
}

func FavouriteGetByVideoUser(favourite *entity.Favourite, uid int64, vid int64) error {
	err := mysqlDB.Where("user_id=?", uid).Where("video_id=?", vid).Find(&favourite).Error
	if err != nil {
		zap.L().Error("mysql:favourites get by vid and uid failed")
		return xerr.ErrDatabase
	}
	return nil
}

func FavouriteUpdate(favourite *entity.Favourite) error {
	err := mysqlDB.Where("user_id=?", favourite.UserID).Where("video_id=?", favourite.VideoID).Update("is_favourite", favourite.IsFavourite).Error
	if err != nil {
		zap.L().Error("mysql:favourites update an entry failed")
		return xerr.ErrDatabase
	}
	return nil
}
