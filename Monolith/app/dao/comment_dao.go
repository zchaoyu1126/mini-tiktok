package dao

import (
	"mini-tiktok/app/entity"
	"mini-tiktok/common/xerr"

	"go.uber.org/zap"
)

func TxCommentCreate(comment *entity.Comment) error {
	err := mysqlDB.Create(comment).Error
	if err != nil {
		zap.L().Error("mysql:comments add new entry failed")
		return xerr.ErrDatabase
	}
	return nil
}

func TxCommentDelete(comment *entity.Comment) error {
	err := mysqlDB.Delete(comment).Error
	if err != nil {
		zap.L().Error("mysql:comments delete an entry failed")
		return xerr.ErrDatabase
	}
	return nil
}

func CommentGetByCommentID(comment *entity.Comment) error {
	err := mysqlDB.Where("comment_id=?", comment.CommentID).First(comment).Error
	if err != nil {
		zap.L().Error("mysql:comments get an entry failed")
		return xerr.ErrDatabase
	}
	return nil
}

func CommentListGetByVideoID(videoID int64, commentList *[]*entity.Comment) error {
	err := mysqlDB.Where("video_id=?", videoID).Find(&commentList).Error
	if err != nil {
		zap.L().Error("mysql:comments get entry list failed")
		return xerr.ErrDatabase
	}
	return nil
}
