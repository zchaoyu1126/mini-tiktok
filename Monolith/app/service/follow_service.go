package service

import (
	"errors"
	"fmt"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"

	"gorm.io/gorm"
)

func RelationAction(fromUID, toUID int64, actionType int64) error {
	if actionType == 1 {
		follow := &entity.Follow{FromUserID: fromUID, ToUserID: toUID, IsFollow: true}
		// 查询记录follow记录是否存在
		err := dao.FollowGetByIDs(follow)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果记录不存在就创建，创建失败则返回nerr，成功则返回nil
				if nerr := dao.FollowAdd(follow); nerr != nil {
					return nerr
				}
				return nil
			} else {
				// 查询失败，返回错误
				return err
			}
		}

		// 记录存在的情况
		follow.IsFollow = true
		if err := dao.FollowUpdate(follow); err != nil {
			return err
		}
		if err := UpdateUserFollowCount(fromUID, 1); err != nil {
			return err
		}
		if err := UpdateUserFollowerCount(toUID, 1); err != nil {
			return err
		}
		return nil
	} else if actionType == 2 {
		// 取消关注，说明数据库中一定是有记录的
		follow := &entity.Follow{FromUserID: fromUID, ToUserID: toUID}
		if err := dao.FollowGetByIDs(follow); err != nil {
			return nil
		}
		follow.IsFollow = false
		if err := dao.FollowUpdate(follow); err != nil {
			return err
		}
		if err := UpdateUserFollowCount(fromUID, -1); err != nil {
			return err
		}
		if err := UpdateUserFollowerCount(toUID, -1); err != nil {
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

// 查询user的关注列表
func FollowList(userID, uid int64) ([]UserVO, error) {
	var follows []entity.Follow
	err := dao.FollowGetByFromUID(userID, &follows)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	userList := make([]UserVO, len(follows))
	for i, follow := range follows {
		user, _ := UserInfo(follow.ToUserID, uid)
		userList[i] = *user
	}
	return userList, nil
}

func FollowerList(userID, uid int64) ([]UserVO, error) {
	var follows []entity.Follow
	err := dao.FollowGetByToUID(userID, &follows)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	userList := make([]UserVO, len(follows))
	fmt.Println(len(follows))
	for i, follow := range follows {
		user, _ := UserInfo(follow.FromUserID, uid)
		userList[i] = *user
	}
	return userList, nil
}
