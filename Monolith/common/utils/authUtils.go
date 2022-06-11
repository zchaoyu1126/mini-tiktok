package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"
)

func GetLoginUser(c *gin.Context) (*entity.User, error) {
	fromUIDTmp, exists := c.Get("uid")
	fromUID, ok := fromUIDTmp.(int64)
	if !exists || !ok {
		zap.S().Errorf("parse fromUserID:%v failed", fromUIDTmp)
		return nil, xerr.ErrBadRequest
	}
	user := &entity.User{UserID: fromUID}
	dao.UserGetByUID(user)
	return user, nil
}
