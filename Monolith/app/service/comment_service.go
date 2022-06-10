package service

import (
	"fmt"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/utils"
	"mini-tiktok/common/xerr"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type CommentVO struct {
	ID         int64  `json:"id"`
	Content    string `json:"content"`
	UserVO     UserVO `json:"user"`
	CreateDate string `json:"create_date"`
}

func CommentAdd(video_id string, comment_text string, user_id int64) (*CommentVO, error) {
	var comment_id int64
	var responseComment *CommentVO
	videoId, _ := strconv.ParseInt(video_id, 10, 64)

	// 生成全局唯一的commit id
	snowFlake, err := utils.NewSnowFlake(2, 2)
	if err != nil {
		zap.L().Error("snow flake init failed")
		return nil, xerr.ErrInternalServer
	}
	comment_id, err = snowFlake.NextId()
	if err != nil {
		zap.L().Error("generate comment id failed")
		return nil, xerr.ErrInternalServer
	}

	comment := &entity.Comment{
		CommentID: comment_id,
		UserID:    user_id,
		VideoID:   videoId,
		Content:   comment_text,
		Status:    1,
	}

	err = dao.CommentAdd(comment)
	if err != nil {
		return nil, err
	}

	// 评论的所有者 userVO
	userVO, err := UserInfo(user_id, user_id)
	if err != nil {
		return nil, err
	}

	responseComment = &CommentVO{
		ID:         comment_id,
		Content:    comment_text,
		UserVO:     *userVO,
		CreateDate: time.Now().Format("mm-dd"),
	}

	return responseComment, nil
}

func CommentDelete(video_id string, comment_id string, user_id int64) (*CommentVO, error) {
	comment := new(entity.Comment)

	commentId, _ := strconv.ParseInt(comment_id, 10, 64)
	comment.CommentID = commentId

	// 首先根据commentID获取这条评论
	err := dao.CommentGetByCommentID(comment)
	if err != nil {
		return nil, err
	}

	// 查询发起人与这条评论作者之间的关系
	userVO, err := UserInfo(comment.UserID, user_id)
	if err != nil {
		return nil, err
	}

	fmt.Println(comment)
	// 由于comment是唯一的，所以根据commentID可直接删除
	err = dao.CommentDelete(comment)
	if err != nil {
		return nil, err
	}

	responseComment := &CommentVO{
		ID:         commentId,
		Content:    comment.Content,
		UserVO:     *userVO,
		CreateDate: comment.CreatedAt.Format("mm-dd"),
	}

	return responseComment, nil
}

func CommentList(user_id int64, video_id string) ([]CommentVO, error) {
	videoId, _ := strconv.ParseInt(video_id, 10, 64)
	var commentList []entity.Comment
	var commentVOList []CommentVO

	// 根据vid得到该视频的所有评论
	err := dao.CommentListGetByVideoId(videoId, &commentList)
	if err != nil {
		return nil, err
	}

	commentVOList = make([]CommentVO, len(commentList))
	// 需要根据commentList[i].UserID去获取UserVO对象
	for i := 0; i < len(commentList); i++ {
		commentVOList[i].ID = commentList[i].CommentID
		commentVOList[i].Content = commentList[i].Content
		zap.S().Infof("%d %d", commentList[i].UserID, user_id)
		userVO, _ := UserInfo(commentList[i].UserID, user_id)
		commentVOList[i].UserVO = *userVO
		commentVOList[i].CreateDate = commentList[i].CreatedAt.Format("mm-dd")
	}

	return commentVOList, err
}
