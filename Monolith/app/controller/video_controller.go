package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []*service.PublicationVO `json:"video_list"`
}

type FeedResponse struct {
	Response
	NextTime  int                      `json:"next_time"`
	VideoList []*service.PublicationVO `json:"video_list"`
}

// 查看用户发布视频列表接口
func PublishList(c *gin.Context) {
	fromUID, err := getUID(c, false)
	if err != nil {
		errorHandler(c, err)
		return
	}
	toUID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	videos, err := service.GetVideoListByUser(toUID, fromUID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		success,
		videos,
	})
}

// 用户发布视频接口
func Publish(c *gin.Context) {
	uid, err := getUID(c, true)
	if err != nil {
		errorHandler(c, err)
		return
	}

	title := c.PostForm("title")
	file, err := c.FormFile("data")
	// 仅支持上传mp4后缀文件
	if err != nil || !strings.HasSuffix(file.Filename, "mp4") {
		errorHandler(c, xerr.ErrInvaildFile)
		return
	}

	videoPath := savePath(file.Filename)
	c.SaveUploadedFile(file, videoPath)

	err = service.VideoPublish(videoPath, title, uid)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, success)
}

// 视频流接口
func Feed(c *gin.Context) {
	uid, err := getUID(c, false)
	if err != nil {
		errorHandler(c, err)
		return
	}
	lastest := c.Query("latest_time")
	if lastest == "" {
		lastest = time.Now().Format("2006-01-02 15:04:05")
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

// 用户上传视频的保存路径为 upload/videos/timestamp_filename
func savePath(fileName string) string {
	var builder strings.Builder
	builder.WriteString("upload/videos/")
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	builder.WriteString(timeStamp)
	builder.WriteString("_" + fileName)
	return builder.String()
}
