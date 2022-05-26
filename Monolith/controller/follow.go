package controller

import (
	"mini-tiktok/service"
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
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}
	actionType, err := strconv.ParseInt(actionTypeTmp, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			UserId:   -1,
			Token:    "",
		})
		return
	}

	service.RelationAction(fromUserID, toUserID, actionType)
	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: ""})
}
