package service

import (
	"errors"
	"gorm.io/gorm"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
)

type FavouriteVideoVo struct {
	PublicationVO
	isFavourite byte
}

func Thumb(uid int64, video_id int64) error {
	return favouriteSave(uid, video_id, 1)
}

func DeThumb(uid int64, video_id int64) error {
	return favouriteSave(uid, video_id, 2)
}

func favouriteSave(uid int64, video_id int64, action_type int) error {
	favourite := &entity.Favourite{
		UserID:  uid,
		VideoID: video_id,
	}
	err := dao.FavouriteGetById(favourite)
	if action_type == 1 {
		favourite.IsFavourite = 1
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				terr := dao.FavouriteAdd(favourite)

				if terr != nil {
					return terr
				}
				updateVideoFavourite(video_id, 1)
			} else {
				return err
			}
		}
		terr := dao.FavouriteUpdate(favourite)
		if terr != nil {
			return terr
		}
		updateVideoFavourite(video_id, 1)
	}
	if action_type == 2 {
		favourite.IsFavourite = 0
		if err != nil {
			return err
		}
		terr := dao.FavouriteUpdate(favourite)
		if terr != nil {
			return terr
		}
		terr = updateVideoFavourite(video_id, -1)
		if terr != nil {
			return terr
		}
	}
	return nil
}

func updateVideoFavourite(vid int64, delta int) error {
	video := &entity.Publication{VideoID: vid}
	err := dao.VideoGetById(video)
	if err != nil {
		return err
	}
	video.FavouriteCount += delta
	err = dao.VideoUpdateById(video)
	if err != nil {
		return err
	}
	return nil
}

func FavouriteVideoList(from_id int64, to_id int64) ([]PublicationVO, error) {
	videos, err := GetVideoListByUser(to_id, from_id)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
