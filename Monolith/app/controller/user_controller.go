package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// 查看用户信息
func UserInfo(c *gin.Context) {
	// fromUID查看toUID的用户信息
	// toUID由query参数中的user_id获得
	// fromUID从userauth中间件的写入的"uid"参数获得，已被封装在getUID函数中
	toUID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fromUID, err := getUID(c, false)
	if err != nil {
		errorHandler(c, xerr.ErrBadRequest)
		return
	}

	user, err := service.UserInfo(toUID, fromUID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		Response: success,
		UserVO:   user,
	})
}
