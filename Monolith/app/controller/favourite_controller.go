package controller

import (
	"mini-tiktok/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavouriteVideoResponse struct {
	response  Response
	VideoList []*service.PublicationVO `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	uid, err := getUID(c, true)
	if err != nil {
		errorHandler(c, err)
		return
	}

	vid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if action == 1 {
		err := service.Thumb(uid, vid)
		if err != nil {
			errorHandler(c, err, "点赞失败")
			return
		}
	} else if action == 2 {
		err := service.DeThumb(uid, vid)
		if err != nil {
			errorHandler(c, err, "取消点赞失败")
			return
		}
	}
	c.JSON(http.StatusOK, success)
}

func FavoriteList(c *gin.Context) {
	fromUID, err := getUID(c, false)
	if err != nil {
		errorHandler(c, err)
		return
	}
	toUID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	videos, err := service.FavouriteVideoList(fromUID, toUID)
	if err != nil {
		errorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, FavouriteVideoResponse{
		response:  success,
		VideoList: videos,
	})
}
