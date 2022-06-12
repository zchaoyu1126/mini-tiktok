package dao

import (
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
)

func VideoAdd(video *entity.Publication) error {
	if err := mysqlDB.Create(video).Error; err != nil {
		zap.L().Error("mysql:publications add new publication failed")
		return xerr.ErrDatabase
	}
	return nil
}

func VideoGetByUser(userID int64, videoList *[]entity.Publication) error {
	err := mysqlDB.Where("owner_id=?", userID).Find(&videoList).Error
	if err != nil {
		zap.L().Error("mysql:publications get video by user failed")
		return xerr.ErrDatabase
	}
	return nil
}

func VideoList(videoList *[]entity.Publication) error {
	err := mysqlDB.Find(&videoList).Error
	if err != nil {
		zap.S().Errorf("mysql:publications get video list failed%w", err)
		return xerr.ErrDatabase
	}
	return nil
}

func VideoGetById(publication *entity.Publication) error {
	err := mysqlDB.Where("video_id", publication.VideoID).Find(publication).Error
	if err != nil {
		zap.S().Errorf("mysql:publications get video by id failed%w", err)
		return xerr.ErrDatabase
	}
	return nil
}

func VideoUpdateById(publication *entity.Publication) error {
	err := mysqlDB.Save(publication).Error
	if err != nil {
		zap.S().Errorf("mysql:publications update failed%w", err)
		return xerr.ErrDatabase
	}
	return nil
}
