package main

import (
	_ "mini-tiktok/common/logger"
	"mini-tiktok/common/utils"
	"mini-tiktok/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {

}
func main() {
	utils.CreateDir("upload/videos/")
	utils.CreateDir("upload/covers/")

	r := gin.New()
	router.InitRouter(r)

	// 开一个文件服务器
	http.Handle("/", http.FileServer(http.Dir("./")))

	go http.ListenAndServe(":8079", nil)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
