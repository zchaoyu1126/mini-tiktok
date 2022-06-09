package dao

import (
	"errors"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UserAdd(user *entity.User) error {
	if err := mysqlDB.Create(user).Error; err != nil {
		zap.L().Error("mysql:users create new user failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 除用户不存在导致的err外，其余所有的error情况都直接返回ErrDatabase
func UserGetByName(user *entity.User) error {
	if err := mysqlDB.Where("username=?", user.UserName).First(user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("mysql:users query user by username failed")
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}

func UserGetByUID(user *entity.User) error {
	if err := mysqlDB.Where("user_id=?", user.UserID).First(user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("mysql:users query user by ID failed")
			return xerr.ErrDatabase
		}
		return err
	}
	return nil
}

func UserUpdate(user *entity.User) error {
	if err := mysqlDB.Save(user).Error; err != nil {
		zap.L().Error("mysql:users update user failed")
		return xerr.ErrDatabase
	}
	return nil
}
