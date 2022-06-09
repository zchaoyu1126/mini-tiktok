package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
