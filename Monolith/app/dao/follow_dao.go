package dao

import (
	"errors"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func FollowAdd(follow *entity.Follow) error {
	err := mysqlDB.Create(follow).Error
	if err != nil {
		zap.L().Error("mysql:follows add new entry failed")
		return xerr.ErrDatabase
	}
	return nil
}

func FollowGetByIDs(follow *entity.Follow) error {
	err := mysqlDB.Where("from_user_id = ? AND to_user_id = ?", follow.FromUserID, follow.ToUserID).First(follow).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("mysql:follows get by ids failed")
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}

func FollowUpdate(follow *entity.Follow) error {
	err := mysqlDB.Save(follow).Error
	if err != nil {
		zap.L().Error("mysql:follows update an entry failed")
		return xerr.ErrDatabase
	}
	return nil
}

func FollowGetByFromUID(fuid int64, follows *[]entity.Follow) error {
	err := mysqlDB.Where("from_user_id = ? AND is_follow = 1", fuid).Find(follows).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("mysql:follows get by fuid failed")
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}

func FollowGetByToUID(tuid int64, follows *[]entity.Follow) error {
	err := mysqlDB.Where("to_user_id = ? AND is_follow = 1", tuid).Find(follows).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("mysql:follows get by tuid failed")
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}
