package dao

import (
	"errors"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 创建一个事务，使用create方法创建favourite记录，同时更新favourite.VideoID的点赞数目
func TxFavouriteAdd(favourite *entity.Favourite) error {
	// 注意在事务中要使用 tx 作为数据库句柄
	// 当panic时，要执行回滚操作
	tx := mysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("mysql:transaction(favourite create) paniced.")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Error("mysql:transaction(favourite create) init failed.")
		return xerr.ErrDatabase
	}

	// 新增点赞记录
	if err := tx.Create(favourite).Error; err != nil {
		zap.L().Error("mysql:transaction(favourite create) create an new entry failed.")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	// 改变视频的点赞数目,首先需要获得原有的数目
	video := &entity.Publication{}
	if err := tx.Where("video_id = ?", favourite.VideoID).First(video).Error; err != nil {
		zap.L().Error("mysql:transaction(favourite create) query video failed.")
		tx.Rollback()
		return xerr.ErrDatabase
	}
	newCount := video.FavouriteCount + 1
	err := tx.Model(video).Where("video_id = ?", favourite.VideoID).Update("favourite_count", newCount).Error
	if err != nil {
		zap.L().Error("mysql:transaction(favourite create) update video failed.")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("mysql:transaction(favourite create) commit failed.")
		return xerr.ErrDatabase
	}
	return nil
}

// 创建一个事务，使用save方法更新favourite记录，同时更新favourite.VideoID的点赞数目
func TxFavouriteUpdate(favourite *entity.Favourite) error {
	tx := mysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("mysql:transaction(favourite create) paniced.")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Error("mysql:transaction(favourite create) init failed.")
		return xerr.ErrDatabase
	}

	// 使用save方法更新favourite
	err := tx.Save(favourite).Error
	if err != nil {
		zap.L().Error("mysql:favourites update an entry failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	video := &entity.Publication{}
	if err := tx.Where("video_id = ?", favourite.VideoID).First(video).Error; err != nil {
		zap.L().Error("mysql:transaction(favourite create) query video failed.")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	var newCount int
	if favourite.IsFavourite {
		newCount = video.FavouriteCount + 1
	} else {
		newCount = video.FavouriteCount - 1
	}

	err = tx.Model(video).Where("video_id = ?", favourite.VideoID).Update("favourite_count", newCount).Error
	if err != nil {
		zap.L().Error("mysql:transaction(favourite create) update video failed.")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	if err := tx.Commit().Error; err != nil {
		zap.L().Error("mysql:transaction(favourite create) commit failed.")
		return xerr.ErrDatabase
	}
	return nil
}

// 根据favourite.UserID和favourite.VideoID 获取点赞记录，有可能返回记录不存在的错误
func FavouriteGetByIDs(favourite *entity.Favourite) error {
	err := mysqlDB.Where("user_id=?", favourite.UserID).Where("video_id=?", favourite.VideoID).First(favourite).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("mysql:favourites get by vid and uid failed")
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}

// 根据userID获取userID所有的点赞记录，不会返回记录不存在错误
func FavouriteGetByUser(userID int64, favs *[]*entity.Favourite) error {
	err := mysqlDB.Where("user_id=?", userID).Find(&favs).Error
	if err != nil {
		zap.L().Error("mysql:favourite get by user failed")
		return xerr.ErrDatabase
	}
	return nil
}
