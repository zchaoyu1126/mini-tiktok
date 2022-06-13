package dao

import (
	"errors"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 创建video
func VideoCreate(video *entity.Publication) error {
	if err := mysqlDB.Create(video).Error; err != nil {
		zap.L().Error("mysql:publications add new publication failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 根据userID获取userID用户发布的所有视频，不会返回记录不存在错误
func VideoGetByUser(userID int64, videoList *[]*entity.Publication) error {
	err := mysqlDB.Where("owner_id=?", userID).Find(&videoList).Error
	if err != nil {
		zap.L().Error("mysql:publications get video by user failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 返回所有的视频
func VideoList(videoList *[]*entity.Publication) error {
	err := mysqlDB.Order("id desc").Find(&videoList).Error
	if err != nil {
		zap.S().Errorf("mysql:publications get video list failed%w", err)
		return xerr.ErrDatabase
	}
	return nil
}

// 根据publication.VideoID获取一条视频，有可能返回记录为空的错误
func VideoGetByID(publication *entity.Publication) error {
	err := mysqlDB.Where("video_id", publication.VideoID).First(publication).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Errorf("mysql:publications get video by id failed%w", err)
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}
