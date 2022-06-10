package dao

import (
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
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
		zap.L().Error("mysql:follows get by ids failed")
		return xerr.ErrDatabase
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
