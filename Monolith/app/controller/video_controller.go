package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"mini-tiktok/app/constant"
	"mini-tiktok/app/entity"
	"mini-tiktok/app/service"
	"mini-tiktok/common/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

type VideoResponse struct {
	Response
	VideoList []service.PublicationVo `json:"video_list"`
}

func PublishList(c *gin.Context) {
	//token := c.Query("token")
	userId := c.Query("user_id")
	uid, err2 := strconv.Atoi(userId)
	if err2 != nil {
		return
	}
	videos, err := service.GetVideoListByUser(int64(uid))
	if err != nil {
		c.JSON(http.StatusOK, VideoResponse{
			Response{StatusCode: 1, StatusMsg: constant.QueryFail},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, VideoResponse{
		Response{StatusCode: 0, StatusMsg: ""},
		videos,
	})
}

// Publish 待解决问题 token鉴权，没有判断上传的是否为视频/*
func Publish(c *gin.Context) {
	file, header, err := c.Request.FormFile("data")
	//token := c.Request.Form.Get("token")
	title := c.Request.Form.Get("title")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  constant.InvalidFile,
		})
		return
	}
	filepath := "videos/" + header.Filename
	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  constant.UnknownError,
		})
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  constant.UnknownError,
		})
		return
	}
	flake, err := utils.NewSnowFlake(1, 1)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  constant.UnknownError,
		})
		return
	}
	vid, _ := flake.NextId()
	now := time.Now()
	// ownerId need to be changed
	err = service.VideoPublish(&entity.Publication{
		VideoID:        vid,
		OwnerID:        1,
		Title:          title,
		PlayUrl:        filepath,
		CoverUrl:       "",
		CreateTime:     &now,
		FavouriteCount: 0,
		CommentCount:   0,
		Status:         0,
	})
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  constant.UnknownError,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0, StatusMsg: "success",
	})
}
