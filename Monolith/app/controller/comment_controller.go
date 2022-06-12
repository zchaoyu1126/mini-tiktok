package controller

import (
	"mini-tiktok/app/service"
	"mini-tiktok/common/xerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	Response
	Comment *service.CommentVO `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []*service.CommentVO `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	//获取参数
	var comment_text string
	var cid int64
	action_type := c.Query("action_type")
	vid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	uid, err := getUID(c, true)
	if err != nil {
		errorHandler(c, err)
		return
	}

	//判断action_type的值，以确定是删除还是增加
	if action_type == "1" {
		//增加
		comment_text = c.Query("comment_text")
		comment, err := service.CommentCreate(vid, comment_text, uid)
		if err != nil {
			errorHandler(c, err)
			return
		}

		c.JSON(
			http.StatusOK,
			CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "add comment success"},
				Comment:  comment,
			},
		)
		return
	} else if action_type == "2" {
		//删除
		cid, _ = strconv.ParseInt(c.Query("comment_id"), 10, 64)
		comment, err := service.CommentDelete(vid, cid, uid)
		if err != nil {
			errorHandler(c, err)
			return
		}
		c.JSON(
			http.StatusOK,
			CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "delete comment success"},
				Comment:  comment,
			},
		)
		return
	}
	errorHandler(c, xerr.ErrBadRequest)
}

func CommentList(c *gin.Context) {
	//获取参数
	uid, err := getUID(c, false)
	if err != nil {
		errorHandler(c, err)
		return
	}
	vid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	commentVoList, err := service.CommentList(uid, vid)
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
