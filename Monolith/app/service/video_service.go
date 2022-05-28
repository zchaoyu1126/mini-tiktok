package service

import (
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
)

type PublicationVo struct {
	ID             int64
	Author         entity.User
	PlayUrl        string
	CoverUrl       string
	FavouriteCount int
	CommentCount   int
	IsFavourite    bool
	title          string
}

func GetVideoListByUser(uid int64) ([]PublicationVo, error) {
	var videos []entity.Publication
	var ret []PublicationVo
	err := dao.VideoGetByUser(uid, &videos)
	if err != nil {
		return nil, err
	}

	for _, v := range videos {
		vo, err := VideoTransform(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, vo)
	}
	return ret, err
}

func VideoPublish(video *entity.Publication) error {
	err := dao.VideoAdd(video)
	if err != nil {
		return err
	}
	return err
}

func VideoTransform(v *entity.Publication) (PublicationVo, error) {
	vo := PublicationVo{
		ID:             v.VideoID,
		PlayUrl:        v.PlayUrl,
		CoverUrl:       v.CoverUrl,
		FavouriteCount: v.FavouriteCount,
		CommentCount:   v.CommentCount,
		title:          v.Title,
	}
	var favourite entity.Favourite
	err := dao.FavouriteGetByVideoUser(&favourite, v.OwnerID, v.VideoID)
	if err != nil {
		return PublicationVo{}, err
	}
	vo.IsFavourite = favourite.IsFavourite == 1
	vo.Author.UserID = v.OwnerID
	err = dao.UserGetByUID(&vo.Author)
	if err != nil {
		return PublicationVo{}, err
	}
	return vo, err
}
