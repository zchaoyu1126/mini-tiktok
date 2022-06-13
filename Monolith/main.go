package main

import (
	_ "mini-tiktok/common/logger"
	"mini-tiktok/common/utils"
	"mini-tiktok/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func appServer() {
	utils.CreateDir("upload/videos/")
	utils.CreateDir("upload/covers/")

	r := gin.New()
	router.InitRouter(r)
	r.Run()
}

func fileServer() {
	http.Handle("/", http.FileServer(http.Dir("./")))

	go http.ListenAndServe(":8079", nil)
}
func main() {

	go fileServer()
	go appServer()

	select {}
}
