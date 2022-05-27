package router

import (
	"mini-tiktok/common/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware(method string) func(c *gin.Context) {
	// https://www.ruanyifeng.com/blog/2019/04/oauth_design.html
	// https://www.ruanyifeng.com/blog/2019/04/oauth-grant-types.html

	return func(c *gin.Context) {
		var token string
		if method == "Header" {
			token = c.PostForm("token")
		} else if method == "Query" {
			token = c.Query("token")
		}
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 1002,
				"msg":  "token为空",
			})
			c.Abort()
			return
		}

		// 去redis中查询token是否过期
		isVaild, err := db.NewRedisDaoInstance().IsTokenValid(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1003,
				"msg":  "redis查询失败",
			})
			c.Abort()
			return
		}

		if !isVaild {
			c.JSON(http.StatusOK, gin.H{
				"code": 1004,
				"msg":  "token已过期",
			})
			c.Abort()
			return
		}

		id, err := db.NewRedisDaoInstance().GetToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1003,
				"msg":  "redis查询失败",
			})
			c.Abort()
			return
		}
		c.Set("fromUserID", id)
		c.Next()
	}
}
