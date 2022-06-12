package service

import (
	"errors"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"gorm.io/gorm"
)

// fromUID 关注 toUID
func Follow(fromUID, toUID int64) error {
	// 查询这条记录follow记录是否存在
	follow := &entity.Follow{FromUserID: fromUID, ToUserID: toUID}
	err := dao.FollowGetByIDs(follow)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 记录不存在时创建，创建失败则返回nerr，成功则返回nil
		follow.IsFollow = true
		terr := dao.TxFollowCreate(follow)
		if terr != nil {
			return terr
		}
		return nil
	} else if err != nil {
		// 查询失败，返回错误
		return err
	}

	// 记录存在的情况
	if follow.IsFollow {
		return xerr.ErrRepeatFollow
	}

	follow.IsFollow = true
	// 使用事务更新关注表，同时增加粉丝数和关注数
	if err := dao.TxFollowUpdate(follow); err != nil {
		return err
	}
	return nil
}

// fromUID 取消关注 UID
func DeFollow(fromUID, toUID int64) error {
	// follow这条记录必定存在
	follow := &entity.Follow{FromUserID: fromUID, ToUserID: toUID}
	if err := dao.FollowGetByIDs(follow); err != nil {
		return nil
	}

	if !follow.IsFollow {
		return xerr.ErrRepeatDeFollow
	}

	follow.IsFollow = false
	// 使用事务更新关注表，同时减少粉丝数和关注数
	if err := dao.TxFollowUpdate(follow); err != nil {
		return err
	}
	return nil
}

// fromUID查看toUID的关注列表
func FollowList(toUID, fromUID int64) ([]*UserVO, error) {
	var follows []*entity.Follow
	if err := dao.FollowGetByFromUID(toUID, &follows); err != nil {
		return nil, err
	}

	userList := make([]*UserVO, len(follows))
	for i, follow := range follows {
		user, _ := UserInfo(follow.ToUserID, fromUID)
		userList[i] = user
	}
	return userList, nil
}

// fromUID查看toUID的粉丝列表
func FollowerList(toUID, fromUID int64) ([]*UserVO, error) {
	var follows []*entity.Follow
	if err := dao.FollowGetByToUID(toUID, &follows); err != nil {
		return nil, err
	}

	userList := make([]*UserVO, len(follows))
	for i, follow := range follows {
		user, _ := UserInfo(follow.FromUserID, fromUID)
		userList[i] = user
	}
	return userList, nil
}
