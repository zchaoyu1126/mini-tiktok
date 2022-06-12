package controller

import (
	"mini-tiktok/common/xerr"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var success Response = Response{0, "success"}

func errorHandler(c *gin.Context, err error, options ...string) {
	var ok bool
	var codeErr *xerr.CodeError

	if codeErr, ok = err.(*xerr.CodeError); !ok {
		codeErr = xerr.ErrUnKnown
	}

	httpCode := codeErr.HTTPCode()
	code := codeErr.ErrCode()
	msg := buildMessage(codeErr.ErrMsg(), options...)
	c.JSON(httpCode, UserLoginResponse{
		Response: Response{code, msg},
	})
}

// 获取由UserAuth中间件写入的uid参数，如果token为0，那么uid参数为-1。
// 当login为true时，说明需要登录，此时如果获取到的uid为-1，那么直接返回未登录错误。
// 当login为false时，不需要登录，允许返回-1，代表游客，允许执行有限的操作。
func getUID(c *gin.Context, login bool) (int64, error) {
	uidTmp, exists := c.Get("uid")
	uid, ok := uidTmp.(int64)
	if !exists || !ok {
		zap.S().Errorf("parse uidTmp:%v failed", uidTmp)
		return 0, xerr.ErrBadRequest
	}
	if uid == -1 && login {
		return -1, xerr.ErrNotLogin
	}
	return uid, nil
}

func buildMessage(msg string, options ...string) string {
	var builder strings.Builder
	builder.WriteString(msg)
	for _, option := range options {
		builder.WriteByte(';')
		builder.WriteString(option)
	}
	return builder.String()
}
