package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/db"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoResponse struct {
	Response
	VideoList []service.PublicationVO `json:"video_list"`
}

type FeedResponse struct {
	Response
	NextTime  int                     `json:"next_time"`
	VideoList []service.PublicationVO `json:"video_list"`
}

func PublishList(c *gin.Context) {
	fromUIDTmp, exists := c.Get("uid")
	fromUID, ok := fromUIDTmp.(int64)
	if !exists || !ok {
		zap.S().Errorf("parse fromUserID:%v failed", fromUIDTmp)
		errorHandler(c, xerr.ErrBadRequest)
		return
	}
	userId := c.Query("user_id")
	uid, err := strconv.Atoi(userId)
	if err != nil {
		zap.S().Errorf("parse user_id:%v failed", uid)
		errorHandler(c, xerr.ErrBadRequest)
		return
	}
	videos, err := service.GetVideoListByUser(int64(uid), fromUID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, VideoResponse{
		success,
		videos,
	})
}

// Publish 待解决问题 token鉴权，没有判断上传的是否为视频/*
func Publish(c *gin.Context) {
	// 获取参数
	file, err := c.FormFile("data")
	title := c.PostForm("title")
	// file, header, err := c.Request.FormFile("data")
	// //token := c.Request.Form.Get("token")
	// title := c.Request.Form.Get("title")
	if err != nil {
		errorHandler(c, xerr.ErrInvaildFile)
		return
	}
	filepath := "upload/videos/" + file.Filename
	c.SaveUploadedFile(file, filepath)

	uidTmp, exists := c.Get("uid")
	uid, ok := uidTmp.(int64)
	if !exists || !ok {
		zap.S().Errorf("parse fromUserID:%v failed", uidTmp)
		errorHandler(c, xerr.ErrBadRequest)
		return
	}

	// out, err := os.Create(filepath)
	// if err != nil {
	// 	zap.L().Error("create video file failed")
	// 	errorHandler(c, xerr.ErrInternalServer)
	// 	return
	// }
	// defer out.Close()

	// _, err = io.Copy(out, file)
	// if err != nil {
	// 	zap.L().Error("copy video file failed")
	// 	errorHandler(c, xerr.ErrInternalServer)
	// 	return
	// }

	// 模拟器的地址
	// url := "http://10.0.2.2:8079/" + filepath
	// 服务器上
	url := "http://39.107.81.188:8079/" + filepath
	err = service.VideoPublish(url, title, uid)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, success)
}

func Feed(c *gin.Context) {
	lastest := c.Query("latest_time")
	if lastest == "" {
		lastest = time.Now().Format("2006-01-02 15:04:05")
	}
	token := c.Query("token")
	var uid int64
	if token == "" {
		uid = -1
	} else {
		uid, _ = db.NewRedisDaoInstance().GetToken(token)
	}

	nextTime, videos, err := service.Feed(lastest, uid)
	if err != nil {
		errorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		success,
		nextTime,
		videos,
	})
}
