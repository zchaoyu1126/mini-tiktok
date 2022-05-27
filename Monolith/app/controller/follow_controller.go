package controller

import (
	"mini-tiktok/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RelationAction(c *gin.Context) {
	toUserIDTmp := c.Query("to_user_id")
	actionTypeTmp := c.Query("action_type")
	fromUserIDTmp, _ := c.Get("fromUserID")
	fromUserID, _ := fromUserIDTmp.(int64)

	toUserID, err := strconv.ParseInt(toUserIDTmp, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	actionType, err := strconv.ParseInt(actionTypeTmp, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	if err := service.RelationAction(fromUserID, toUserID, actionType); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "success"})
}
