package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []*service.UserVO `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	toUID, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	fromUID, err := getUID(c, true)
	if err != nil {
		errorHandler(c, err)
		return
	}

	if fromUID == toUID {
		errorHandler(c, xerr.ErrFollowMyself)
		return
	}

	if actionType == 1 {
		err = service.Follow(fromUID, toUID)
	} else if actionType == 2 {
		err = service.DeFollow(fromUID, toUID)
	}
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, success)
}

func FollowList(c *gin.Context) {
	toUID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fromUID, err := getUID(c, false)
	if err != nil {
		errorHandler(c, err)
		return
	}

	list, err := service.FollowList(toUID, fromUID)
	if err != nil {
		errorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		success,
		list,
	})
}

func FollowerList(c *gin.Context) {
	toUID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fromUID, err := getUID(c, false)
	if err != nil {
		errorHandler(c, err)
		return
	}

	list, err := service.FollowerList(toUID, fromUID)
	if err != nil {
		errorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		success,
		list,
	})
}
