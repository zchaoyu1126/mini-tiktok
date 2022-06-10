package controller

import (
	"mini-tiktok/common/xerr"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var success Response = Response{0, "success"}

func errorHandler(c *gin.Context, err error) {
	var ok bool
	var codeErr *xerr.CodeError

	if codeErr, ok = err.(*xerr.CodeError); !ok {
		codeErr = xerr.ErrUnKnown
	}

	httpCode := codeErr.HTTPCode()
	code := codeErr.ErrCode()
	msg := codeErr.ErrMsg()
	c.JSON(httpCode, UserLoginResponse{
		Response: Response{code, msg},
	})
}
