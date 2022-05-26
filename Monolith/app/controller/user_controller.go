package controller

import (
	"mini-tiktok/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserVO struct {
	ID            int64  `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	UserVO `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "because of database, login failed"},
			UserId:   -1,
			Token:    "",
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		UserId:   id,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "login success"},
		UserId:   id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	toUserIDTmp := c.Query("user_id")
	toUserID, _ := strconv.ParseInt(toUserIDTmp, 10, 64)

	fromUserIDTmp, _ := c.Get("fromUserID")
	fromUserID, _ := fromUserIDTmp.(int64)

	userDTO, err := service.UserInfo(toUserID, fromUserID)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	user := UserVO{
		ID:            userDTO.ID,
		UserName:      userDTO.UserName,
		FollowCount:   userDTO.FollowCount,
		FollowerCount: userDTO.FollowerCount,
		IsFollow:      userDTO.IsFollow,
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		UserVO:   user,
	})
}
