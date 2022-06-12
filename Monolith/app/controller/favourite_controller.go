package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mini-tiktok/app/service"
	"mini-tiktok/common/utils"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"
)

type FavouriteVideoResponse struct {
	response  Response
	VideoList []service.PublicationVO `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	user, err2 := utils.GetLoginUser(c)
	if err2 != nil {
		errorHandler(c, err2)
		return
	}
	action := c.Query("action_type")
	//uid, _ := strconv.Atoi(c.Query("user_id"))
	target, _ := strconv.Atoi(c.Query("video_id"))
	ac, _ := strconv.Atoi(action)
	if ac == 1 {
		err := service.Thumb(user.UserID, int64(target))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
			return
		}
	} else if ac == 2 {
		err := service.DeThumb(user.UserID, int64(target))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消点赞失败"})
			return
		}
	}
	c.JSON(http.StatusOK, success)
}

func FavoriteList(c *gin.Context) {
	user, err2 := utils.GetLoginUser(c)
	if err2 != nil {
		errorHandler(c, err2)
		return
	}
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		zap.S().Errorf("parse user_id:%v failed", uid)
		errorHandler(c, xerr.ErrBadRequest)
		return
	}
	videos, err := service.FavouriteVideoList(user.UserID, int64(uid))
	if err != nil {
		zap.S().Errorf("get favourite video list:%v failed", uid)
		errorHandler(c, xerr.ErrDatabase)
		return
	}
	c.JSON(http.StatusOK, FavouriteVideoResponse{
		response:  success,
		VideoList: videos,
	})
}
