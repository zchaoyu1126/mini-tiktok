package router

import (
	"mini-tiktok/app/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	// r.Use(ZapLogger(zap.L()))
	// r.Use(ZapRecovery(zap.L(), true))

	r.StaticFS("/upload/videos", http.Dir("upload/videos"))

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", UserAuth("Query"), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", UserAuth("Body"), controller.Publish)
	apiRouter.GET("/publish/list/", UserAuth("Query"), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", UserAuth("Query"), controller.CommentAction)
	apiRouter.GET("/comment/list/", UserAuth("Query"), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", UserAuth("Query"), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", UserAuth("Query"), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", UserAuth("Query"), controller.FollowerList)
}
