package controller

import (
	"fmt"
	"mini-tiktok/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 需要做一个缓存来记录用户是否存在，用redis缓存

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

	// should redis check wheather this user has already existed
	exist, err := service.CheckUserExist(username)
	if err != nil {
		// 处理service层传递上来的错误
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
	exist, err := service.CheckUserExist(username)
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
	data, _ := c.Get("id")
	id, ok := data.(int64)
	fmt.Println(data, id, ok)
	data, _ = c.Get("callerid")
	callerID, _ := data.(int64)

	userInfomation, err := service.UserInfo(id, callerID)
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
		FollowerCount: userInfomation.FollowCount,
		IsFollow:      userInfomation.IsFollow,
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		User:     user,
	})
}
