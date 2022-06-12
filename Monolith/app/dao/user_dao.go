package dao

import (
	"errors"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 新建一个用户
func UserCreate(user *entity.User) error {
	if err := mysqlDB.Create(user).Error; err != nil {
		zap.L().Error("mysql:users create new user failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 根据用户的username获取用户记录，除用户不存在导致的err外，其余所有的error情况都直接返回ErrDatabase
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

// 根据用户的userid获取用户记录，除用户不存在导致的err外，其余所有的error情况都直接返回ErrDatabase
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
