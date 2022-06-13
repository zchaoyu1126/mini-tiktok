package service

import (
	"errors"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/utils"
	"mini-tiktok/common/xerr"
	"strconv"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
)

type PublicationVO struct {
	ID             int64   `json:"id"`
	Author         *UserVO `json:"author"`
	PlayUrl        string  `json:"play_url"`
	CoverUrl       string  `json:"cover_url"`
	FavouriteCount int     `json:"favourite_count"`
	CommentCount   int     `json:"comment_count"`
	IsFavourite    bool    `json:"is_favourite"`
	Title          string  `json:"titile"`
}

// fromUID查看toUID的视频列表
func GetVideoListByUser(toUID, fromUID int64) ([]*PublicationVO, error) {
	var videos []*entity.Publication
	var ret []*PublicationVO

	// 获取视频列表
	if err := dao.VideoGetByUser(toUID, &videos); err != nil {
		return nil, err
	}

	// 组装视频VO，需要加上fromUID是否点赞了该视频，以及该视频的作者信息
	ret = make([]*PublicationVO, len(videos))
	for i, v := range videos {
		vo, err := VideoTransform(v, fromUID)
		if err != nil {
			return nil, err
		}
		ret[i] = vo
	}
	return ret, nil
}

func VideoPublish(filePath, title string, uid int64) error {
	//host := "http://10.0.2.2:8079/"
	host := "http://39.107.81.188:8079/"

	flake, err := utils.NewSnowFlake(1, 1)
	if err != nil {
		return xerr.ErrInternalServer
	}
	vid, _ := flake.NextId()

	coverPath := savePath()
	extractFrame(filePath, coverPath)

	err = dao.VideoCreate(&entity.Publication{
		VideoID:  vid,
		OwnerID:  uid,
		Title:    title,
		PlayUrl:  host + filePath,
		CoverUrl: host + coverPath,
	})

	if err != nil {
		return err
	}
	return nil
}

func VideoTransform(v *entity.Publication, uid int64) (*PublicationVO, error) {
	vo := &PublicationVO{
		ID:             v.VideoID,
		PlayUrl:        v.PlayUrl,
		CoverUrl:       v.CoverUrl,
		FavouriteCount: v.FavouriteCount,
		CommentCount:   v.CommentCount,
		Title:          v.Title,
	}

	// 查询作者信息
	author, err := UserInfo(v.OwnerID, uid)
	if err != nil {
		return nil, err
	}
	vo.Author = author

	res, err := isFavourite(uid, v.VideoID)
	if err != nil {
		return nil, err
	}
	vo.IsFavourite = res
	return vo, nil
}

func Feed(lastest string, uid int64) (int, []*PublicationVO, error) {
	var videos []*entity.Publication
	var ret []*PublicationVO

	if err := dao.VideoList(&videos); err != nil {
		return 0, nil, err
	}
	ret = make([]*PublicationVO, len(videos))
	for i, v := range videos {
		author, _ := UserInfo(v.OwnerID, uid)
		vo, err := VideoTransform(v, uid)
		if err != nil {
			return 0, nil, err
		}
		vo.Author = author
		ret[i] = vo
	}

	return 0, ret, nil
}

// 查看uid是否点赞了vid视频
func isFavourite(uid, vid int64) (bool, error) {
	if uid == -1 {
		return false, nil
	}
	favourite := &entity.Favourite{UserID: uid, VideoID: vid}
	err := dao.FavouriteGetByIDs(favourite)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return favourite.IsFavourite, nil
}

func extractFrame(videoPath string, framePath string) error {
	err := ffmpeg.Input(videoPath).
		Output(framePath, ffmpeg.KwArgs{
			"vframes": "20",
		}).
		OverWriteOutput().
		Run()
	return err
}

func savePath() string {
	var builder strings.Builder
	builder.WriteString("upload/videos/")
	builder.WriteString(strconv.Itoa(int(time.Now().Unix())))
	builder.WriteString(".jpg")
	return builder.String()
}
