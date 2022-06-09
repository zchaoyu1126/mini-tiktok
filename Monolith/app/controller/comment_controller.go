package controller

import (
	"log"
	"mini-tiktok/app/service"
	"mini-tiktok/common/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type Comment struct {
	ID         int64  `json:"id"`
	Content    string `json:"content"`
	UserVO     UserVO `json:"user"`
	CreateDate string `json:"create_date"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	//获取参数

	var comment_text string
	var comment_id string
	action_type := c.Query("action_type")
	video_id := c.Query("video_id")
	token := c.Query("token")

	//验证token
	user_id, err := db.NewRedisDaoInstance().GetToken(token)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			CommentResponse{
				Response: Response{StatusCode: -1, StatusMsg: "invaild token"},
				Comment:  Comment{},
			},
		)
		return
	}

	//判断action_type的值，以确定是删除还是增加
	if action_type == "1" {
		//增加
		comment_text = c.Query("comment_text")
		sComment, err := service.CommentAdd(video_id, comment_text, user_id)
		if err != nil {
			log.Panicln(err)
			return
		}
		c.JSON(
			http.StatusBadRequest,
			CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "add comment success"},
				Comment: Comment{
					ID:         sComment.ID,
					Content:    sComment.Content,
					UserVO:     UserVO(sComment.UserDTO),
					CreateDate: sComment.CreateDate,
				},
			},
		)
		return
	} else if action_type == "2" {
		//删除
		comment_id = c.Query("comment_text")
		sComment, err := service.CommentDelete(video_id, comment_id, user_id)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(
			http.StatusOK,
			CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "delete comment success"},
				Comment: Comment{
					ID:         sComment.ID,
					Content:    sComment.Content,
					UserVO:     UserVO(sComment.UserDTO),
					CreateDate: sComment.CreateDate,
				},
			},
		)
		return
	} else {
		//action_type 参数错误
		c.JSON(
			http.StatusBadRequest,
			CommentResponse{
				Response: Response{StatusCode: -1, StatusMsg: "the value of action_type is wrong"},
				Comment:  Comment{},
			},
		)
		return
	}

}

func CommentList(c *gin.Context) {
	//获取参数
	token := c.Query("token")
	videoId := c.Query("video_id")
	var commentList []Comment

	//验证token
	user_id, err := db.NewRedisDaoInstance().GetToken(token)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			CommentListResponse{
				Response:    Response{StatusCode: -1, StatusMsg: "invaild token"},
				CommentList: commentList,
			},
		)
		return
	}

	commentVoList, err1 := service.CommentList(user_id, videoId)
	if err1 != nil {
		log.Println(err)
		return
	}

	for i := 0; i < len(commentVoList); i++ {
		commentList[i].ID = commentVoList[i].ID
		commentList[i].Content = commentVoList[i].Content
		commentList[i].UserVO = UserVO(commentVoList[i].UserDTO)
		commentList[i].CreateDate = commentVoList[i].CreateDate
	}

	c.JSON(
		http.StatusOK,
		CommentListResponse{
			Response:    Response{StatusCode: 0, StatusMsg: "SUCCESS"},
			CommentList: commentList,
		},
	)
	return

}
