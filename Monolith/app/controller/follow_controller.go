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
	UserList []service.UserVO `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	toUIDTmp := c.Query("to_user_id")
	actionTypeTmp := c.Query("action_type")
	fromUIDTmp, _ := c.Get("uid")
	fromUID, _ := fromUIDTmp.(int64)

	toUID, err := strconv.ParseInt(toUIDTmp, 10, 64)
	if err != nil {
		errorHandler(c, xerr.ErrBadRequest)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeTmp, 10, 64)
	if err != nil {
		errorHandler(c, xerr.ErrBadRequest)
		return
	}

	err = service.RelationAction(fromUID, toUID, actionType)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, success)
}

func FollowList(c *gin.Context) {
	userIDTmp := c.Query("user_id")
	uidTmp, _ := c.Get("uid")
	userID, _ := strconv.ParseInt(userIDTmp, 10, 64)
	uid, _ := uidTmp.(int64)

	list, err := service.FollowList(userID, uid)
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
	userIDTmp := c.Query("user_id")
	uidTmp, _ := c.Get("uid")
	userID, _ := strconv.ParseInt(userIDTmp, 10, 64)
	uid, _ := uidTmp.(int64)

	list, err := service.FollowerList(userID, uid)
	if err != nil {
		errorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		success,
		list,
	})
}
