package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	Response
	Comment service.CommentVO `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []service.CommentVO `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	//获取参数
	var comment_text string
	var comment_id string
	action_type := c.Query("action_type")
	video_id := c.Query("video_id")
	uidTmp, _ := c.Get("uid")
	uid, _ := uidTmp.(int64)

	//判断action_type的值，以确定是删除还是增加
	if action_type == "1" {
		//增加
		comment_text = c.Query("comment_text")
		comment, err := service.CommentAdd(video_id, comment_text, uid)
		if err != nil {
			errorHandler(c, err)
			return
		}

		c.JSON(
			http.StatusOK,
			CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "add comment success"},
				Comment:  *comment,
			},
		)
		return
	} else if action_type == "2" {
		//删除
		comment_id = c.Query("comment_id")
		comment, err := service.CommentDelete(video_id, comment_id, uid)
		if err != nil {
			errorHandler(c, err)
			return
		}
		c.JSON(
			http.StatusOK,
			CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "delete comment success"},
				Comment:  *comment,
			},
		)
		return
	}
	errorHandler(c, xerr.ErrBadRequest)
}

func CommentList(c *gin.Context) {
	//获取参数
	uidTmp, _ := c.Get("uid")
	uid, _ := uidTmp.(int64)
	videoId := c.Query("video_id")

	commentVoList, err := service.CommentList(uid, videoId)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		CommentListResponse{
			Response:    Response{StatusCode: 0, StatusMsg: "SUCCESS"},
			CommentList: commentVoList,
		},
	)
}
