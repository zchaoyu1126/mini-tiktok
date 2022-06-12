package service

import (
	"errors"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"gorm.io/gorm"
)

// uid点赞vid
func Thumb(uid int64, vid int64) error {
	favourite := &entity.Favourite{
		UserID:  uid,
		VideoID: vid,
	}
	err := dao.FavouriteGetByIDs(favourite)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 记录不存在，使用事务新增点赞记录，同时修改Video的点赞数目
		favourite.IsFavourite = true
		terr := dao.TxFavouriteAdd(favourite)
		if terr != nil {
			return terr
		}
		return nil
	} else if err != nil {
		return err
	}

	// 记录已经存在，首先判断是否为重复点赞，若不是则更新。
	if favourite.IsFavourite {
		return xerr.ErrRepeatThumb
	}
	favourite.IsFavourite = true
	err = dao.TxFavouriteUpdate(favourite)
	if err != nil {
		return err
	}
	return nil
}

// uid取消点赞vid
func DeThumb(uid int64, vid int64) error {
	favourite := &entity.Favourite{
		UserID:  uid,
		VideoID: vid,
	}

	if !favourite.IsFavourite {
		return xerr.ErrRepeatDeThumb
	}
	favourite.IsFavourite = false
	err := dao.TxFavouriteUpdate(favourite)
	if err != nil {
		return err
	}
	return nil
}

// fromUID查看toUID的点赞列表
func FavouriteVideoList(fromUID int64, toUID int64) ([]*PublicationVO, error) {
	var entries []*entity.Favourite
	var ret []*PublicationVO

	// 获取toUID的所有点赞记录
	err := dao.FavouriteGetByUser(toUID, &entries)
	if err != nil {
		return nil, err
	}

	ret = make([]*PublicationVO, len(entries))
	for i, entry := range entries {
		// 根据点赞记录获取视频
		video := &entity.Publication{VideoID: entry.VideoID}
		if err := dao.VideoGetByID(video); err != nil {
			return nil, err
		}

		vo, err := VideoTransform(video, fromUID)
		if err != nil {
			return nil, err
		}
		ret[i] = vo
	}
	return ret, nil
}
