package service

import (
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/utils"
	"mini-tiktok/common/xerr"
)

type PublicationVO struct {
	ID             int64   `json:"id"`
	Author         *UserVO `json:"author"`
	PlayUrl        string  `json:"play_url"`
	CoverUrl       string
	FavouriteCount int
	CommentCount   int
	IsFavourite    bool
	Title          string
}

func GetVideoListByUser(toUID, fromUID int64) ([]PublicationVO, error) {
	var videos []entity.Publication
	var ret []PublicationVO
	// 查询是否关注了该作者
	isFollow := false
	if toUID != fromUID {
		follow := &entity.Follow{FromUserID: fromUID, ToUserID: toUID}
		if err := dao.FollowGetByIDs(follow); err != nil {
			return nil, err
		}
		isFollow = follow.IsFollow
	}

	// 查询作者信息
	owner := &entity.User{UserID: fromUID}
	if err := dao.UserGetByUID(owner); err != nil {
		return nil, err
	}
	// 组装UserVO
	author := &UserVO{owner.ID, owner.UserName, owner.FollowCount, owner.FollowerCount, isFollow}

	if err := dao.VideoGetByUser(toUID, &videos); err != nil {
		return nil, err
	}

	for _, v := range videos {
		vo, err := VideoTransform(&v)
		if err != nil {
			return nil, err
		}
		vo.Author = author
		ret = append(ret, vo)
	}
	return ret, nil
}

func VideoPublish(filepath, title string, uid int64) error {
	flake, err := utils.NewSnowFlake(1, 1)
	if err != nil {
		return xerr.ErrInternalServer
	}
	vid, _ := flake.NextId()

	// 需要生成封面
	err = dao.VideoAdd(&entity.Publication{
		VideoID:  vid,
		OwnerID:  uid,
		Title:    title,
		PlayUrl:  filepath,
		CoverUrl: "",
	})

	if err != nil {
		return err
	}
	return nil
}

func VideoTransform(v *entity.Publication) (PublicationVO, error) {
	vo := PublicationVO{
		ID:             v.VideoID,
		PlayUrl:        v.PlayUrl,
		CoverUrl:       v.CoverUrl,
		FavouriteCount: v.FavouriteCount,
		CommentCount:   v.CommentCount,
		Title:          v.Title,
	}
	var favourite entity.Favourite
	err := dao.FavouriteGetByVideoUser(&favourite, v.OwnerID, v.VideoID)
	if err != nil {
		return PublicationVO{}, err
	}
	vo.IsFavourite = favourite.IsFavourite == 1
	return vo, err
}

func Feed(lastest string, uid int64) (int, []PublicationVO, error) {
	var videos []entity.Publication
	var ret []PublicationVO
	if err := dao.VideoList(&videos); err != nil {
		return 0, []PublicationVO{}, err
	}

	for _, v := range videos {
		if uid == -1 {
			uid = v.OwnerID
		}
		author, _ := UserInfo(v.OwnerID, uid)
		ret = append(ret, PublicationVO{
			ID:             v.VideoID,
			PlayUrl:        v.PlayUrl,
			CoverUrl:       v.CoverUrl,
			FavouriteCount: v.FavouriteCount,
			CommentCount:   v.CommentCount,
			Title:          v.Title,
			Author:         author,
		})
	}

	return 0, ret, nil
}
