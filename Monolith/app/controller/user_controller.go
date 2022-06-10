package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserLoginResponse struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	*service.UserVO `json:"user"`
}

// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.Register(username, password)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: success,
		UserID:   id,
		Token:    token,
	})
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.Login(username, password)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: success,
		UserID:   id,
		Token:    token,
	})
}

// 用户信息
func UserInfo(c *gin.Context) {
	// 获取Query参数 user_id, user_id解析失败时返回ErrBadRequest错误
	toUserIDTmp := c.Query("user_id")
	toUserID, err := strconv.ParseInt(toUserIDTmp, 10, 64)
	if err != nil {
		zap.S().Errorf("parse user_id:%v failed", toUserIDTmp)
		errorHandler(c, xerr.ErrBadRequest)
	}

	// 获取经过路由中间件UserAuthMiddlerWare解析后又写入的fromUserID信息
	fromUIDTmp, exists := c.Get("uid")
	fromUID, ok := fromUIDTmp.(int64)
	if !exists || !ok {
		zap.S().Errorf("parse fromUserID:%v failed", fromUIDTmp)
		errorHandler(c, xerr.ErrBadRequest)
		return
	}

	user, err := service.UserInfo(toUserID, fromUID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		Response: success,
		UserVO:   user,
	})
}
