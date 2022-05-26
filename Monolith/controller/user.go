package controller

import (
	"mini-tiktok/repository"
	"mini-tiktok/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
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
	User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// query redis to check wheather this user has already existed or not
	exist, err := repository.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}
	if exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
			UserId:   -1,
			Token:    "",
		})
		return
	}
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

	// should redis check wheather this user has already existed
	exist, err := repository.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}
	if !exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			UserId:   -1,
			Token:    "",
		})
		return
	}

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

	userInfomation, err := service.UserInfo(toUserID, fromUserID)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	user := User{
		ID:            userInfomation.ID,
		UserName:      userInfomation.UserName,
		FollowCount:   userInfomation.FollowCount,
		FollowerCount: userInfomation.FollowerCount,
		IsFollow:      userInfomation.IsFollow,
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		User:     user,
	})
}
