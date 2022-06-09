package router

import (
	"fmt"
	"mini-tiktok/common/db"
	"mini-tiktok/common/xerr"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func errorHandler(c *gin.Context, err error) {
	var ok bool
	var codeErr *xerr.CodeError

	if codeErr, ok = err.(*xerr.CodeError); !ok {
		codeErr = xerr.ErrUnKnown
	}

	httpCode := codeErr.HTTPCode()
	code := codeErr.ErrCode()
	msg := codeErr.ErrMsg()
	c.JSON(httpCode, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func UserAuth(method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if method == "Body" {
			token = c.PostForm("token")
		} else if method == "Query" {
			token = c.Query("token")
		}
		fmt.Println("hihi", token)
		// 未携带有效的token
		if token == "" {
			errorHandler(c, xerr.ErrTokenValidation)
			c.Abort()
			return
		}

		// 去redis中查询token是否过期
		isVaild, err := db.NewRedisDaoInstance().IsTokenValid(token)
		if err != nil {
			errorHandler(c, err)
			c.Abort()
			return
		}

		// redis已过期
		if !isVaild {
			errorHandler(c, xerr.ErrTokenNotFound)
			c.Abort()
			return
		}

		// 根据token获取userid
		id, err := db.NewRedisDaoInstance().GetToken(token)
		if err != nil {
			errorHandler(c, err)
			c.Abort()
			return
		}
		c.Set("uid", id)
		c.Next()
	}
}

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func ZapRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
