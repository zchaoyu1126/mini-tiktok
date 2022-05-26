package service

import (
	"errors"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"

	"gorm.io/gorm"
)

func RelationAction(fromUserID, toUserID int64, actionType int64) error {
	if actionType == 1 {
		follow := &entity.Follow{FromUserID: fromUserID, ToUserID: toUserID, IsFollow: true}
		// record already exist
		if err := dao.FollowGetByIDs(follow); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := dao.FollowCreate(follow); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		// update record
		if err := dao.FollowUpdate(follow); err != nil {
			return err
		}
		if err := UpdateUserFollowCount(fromUserID, 1); err != nil {
			return err
		}
		if err := UpdateUserFollowerCount(toUserID, 1); err != nil {
			return err
		}
		return nil
	} else if actionType == 2 {
		// update record
		follow := &entity.Follow{FromUserID: fromUserID, ToUserID: toUserID, IsFollow: false}
		if err := dao.FollowUpdate(follow); err != nil {
			return err
		}
		if err := UpdateUserFollowCount(fromUserID, -1); err != nil {
			return err
		}
		if err := UpdateUserFollowerCount(toUserID, -1); err != nil {
			return err
		}
	}
	return nil
}

func UpdateUserFollowCount(userID, delta int64) error {
	user := &entity.User{UserID: userID}
	if err := dao.UserGetByUID(user); err != nil {
		return err
	}
	user.FollowCount += delta
	if err := dao.UserUpdate(user); err != nil {
		return err
	}
	return nil
}

func UpdateUserFollowerCount(userID, delta int64) error {
	user := &entity.User{UserID: userID}
	if err := dao.UserGetByUID(user); err != nil {
		return err
	}
	user.FollowerCount += delta
	if err := dao.UserUpdate(user); err != nil {
		return err
	}
	return nil
}
