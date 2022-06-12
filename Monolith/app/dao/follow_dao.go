package dao

import (
	"errors"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 创建一个事务，使用create方法创建follow记录，并更新follow.FromUserID和ToUserID的关注数或粉丝数
func TxFollowCreate(follow *entity.Follow) error {
	tx := mysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("mysql: transaction(follow create) paniced.")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Error("mysql: transaction(follow create) init failed.")
		return xerr.ErrDatabase
	}

	err := tx.Create(follow).Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow create) add new entry failed")
		return xerr.ErrDatabase
	}

	fromUser := &entity.User{}
	if err := tx.Where("user_id = ? ", follow.FromUserID).First(fromUser).Error; err != nil {
		zap.L().Error("mysql:transaction(follow create) query from user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}
	newCount := fromUser.FollowCount + 1
	err = tx.Model(fromUser).Where("user_id = ?", follow.FromUserID).Update("follow_count", newCount).Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow create) update from user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	toUser := &entity.User{}
	if err := tx.Where("user_id = ? ", follow.ToUserID).First(toUser).Error; err != nil {
		zap.L().Error("mysql:transaction(follow create) query to user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}
	newCount = fromUser.FollowerCount + 1
	err = tx.Model(toUser).Where("user_id = ?", follow.ToUserID).Update("follower_count", newCount).Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow create) update to user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	err = tx.Commit().Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow create) commit failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 创建一个事务，使用save方法更新follow记录，并更新follow.FromUserID和ToUserID的关注数或粉丝数
func TxFollowUpdate(follow *entity.Follow) error {
	tx := mysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("mysql: transaction(follow update) paniced.")
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Error("mysql: transaction(follow update) init failed.")
		return xerr.ErrDatabase
	}

	err := tx.Save(follow).Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow update) add new entry failed")
		return xerr.ErrDatabase
	}

	var delta int64
	if follow.IsFollow {
		delta = 1
	} else {
		delta = -1
	}

	fromUser := &entity.User{}
	if err := tx.Where("user_id = ? ", follow.FromUserID).First(fromUser).Error; err != nil {
		zap.L().Error("mysql:transaction(follow update) query from user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}
	newCount := fromUser.FollowCount + delta
	err = tx.Model(fromUser).Where("user_id = ?", follow.FromUserID).Update("follow_count", newCount).Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow update) update from user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	toUser := &entity.User{}
	if err := tx.Where("user_id = ? ", follow.ToUserID).First(fromUser).Error; err != nil {
		zap.L().Error("mysql:transaction(follow update) query to user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}
	newCount = fromUser.FollowerCount + delta
	err = tx.Model(toUser).Where("user_id = ?", follow.ToUserID).Update("follower_count", newCount).Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow update) update to user failed")
		tx.Rollback()
		return xerr.ErrDatabase
	}

	err = tx.Commit().Error
	if err != nil {
		zap.L().Error("mysql:transaction(follow update) commit failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 根据follow.FromUserID和ToUserID获取数据库中的记录，有可能返回不存在的错误
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

// 根据fuid获取fuid所有的关注者，不会返回记录不存在错误
func FollowGetByFromUID(fuid int64, follows *[]*entity.Follow) error {
	err := mysqlDB.Where("from_user_id = ? AND is_follow = 1", fuid).Find(&follows).Error
	if err != nil {
		zap.L().Error("mysql:follows get by fuid failed")
		return xerr.ErrDatabase
	}
	return nil
}

// 根据tuid获取tuid的所有粉丝，不会返回记录不存在错误
func FollowGetByToUID(tuid int64, follows *[]*entity.Follow) error {
	err := mysqlDB.Where("to_user_id = ? AND is_follow = 1", tuid).Find(&follows).Error
	if err != nil {
		zap.L().Error("mysql:follows get by tuid failed")
		return xerr.ErrDatabase
	}
	return nil
}
