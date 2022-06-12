package service

import (
	"fmt"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/utils"
	"mini-tiktok/common/xerr"
	"time"

	"go.uber.org/zap"
)

type CommentVO struct {
	ID         int64   `json:"id"`
	Content    string  `json:"content"`
	UserVO     *UserVO `json:"user"`
	CreateDate string  `json:"create_date"`
}

// uid给vid添加内容为comment_text的评论
func CommentCreate(vid int64, comment_text string, uid int64) (*CommentVO, error) {
	var comment_id int64
	var responseComment *CommentVO

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
		UserID:    uid,
		VideoID:   vid,
		Content:   comment_text,
	}

	err = dao.TxCommentCreate(comment)
	if err != nil {
		return nil, err
	}

	// 评论的所有者 userVO
	userVO, err := UserInfo(uid, uid)
	if err != nil {
		return nil, err
	}

	responseComment = &CommentVO{
		ID:         comment_id,
		Content:    comment_text,
		UserVO:     userVO,
		CreateDate: time.Now().Format("mm-dd"),
	}

	return responseComment, nil
}

// uid删除vid视频下编号为cid的评论
func CommentDelete(vid int64, cid int64, uid int64) (*CommentVO, error) {
	comment := new(entity.Comment)
	comment.CommentID = cid

	// 首先根据commentID获取这条评论
	err := dao.CommentGetByCommentID(comment)
	if err != nil {
		return nil, err
	}

	// 查询删除请求发起人与这条评论作者之间的关系
	userVO, err := UserInfo(comment.UserID, uid)
	if err != nil {
		return nil, err
	}

	// 由于comment是唯一的，所以根据commentID可直接删除
	err = dao.TxCommentDelete(comment)
	if err != nil {
		return nil, err
	}

	commentVO := &CommentVO{
		ID:         comment.CommentID,
		Content:    comment.Content,
		UserVO:     userVO,
		CreateDate: comment.CreatedAt.Format("mm-dd"),
	}

	return commentVO, nil
}

func CommentList(uid int64, vid int64) ([]*CommentVO, error) {
	var commentList []*entity.Comment
	// 根据vid得到该视频的所有评论
	err := dao.CommentListGetByVideoID(vid, &commentList)
	if err != nil {
		return nil, err
	}

	commentVOList := make([]*CommentVO, len(commentList))
	fmt.Println(len(commentVOList), len(commentList))
	// 需要根据commentList[i].UserID去获取UserVO对象
	for i := 0; i < len(commentList); i++ {
		userVO, _ := UserInfo(commentList[i].UserID, uid)
		commentVOList[i] = &CommentVO{
			ID:         commentList[i].CommentID,
			Content:    commentList[i].Content,
			UserVO:     userVO,
			CreateDate: commentList[i].CreatedAt.Format("mm-dd"),
		}

	}

	return commentVOList, err
}
