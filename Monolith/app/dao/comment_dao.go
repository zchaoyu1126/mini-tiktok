package dao

import (
	"log"
	"mini-tiktok/app/entity"
)

func CommentAdd(comment *entity.Comment) error {
	mysqlDB.AutoMigrate(&entity.Comment{})
	err := mysqlDB.Create(comment).Error
	if err != nil {
		log.Println(err)
	}

	return err
}

func CommentDelete(comment *entity.Comment) error {
	mysqlDB.AutoMigrate(&entity.Comment{})
	err := mysqlDB.Where("comment_id=?", comment.CommentID).Update("status", 0).Error
	if err != nil {
		log.Println(err)
	}

	return err
}

func CommentGetByCommentID(comment *entity.Comment) error {
	mysqlDB.AutoMigrate(&entity.Comment{})
	err := mysqlDB.Where("comment_id=?", comment.CommentID).First(comment).Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func CommentListGetByVideoId(userID int64, commentList *[]entity.Comment) error {
	mysqlDB.AutoMigrate(&entity.Comment{})
	err := mysqlDB.Where("user_id=?", userID).Find(&commentList).Error
	if err != nil {
		log.Println(err)
	}

	return err
}

// a dao used in comment service
func UserGetByVideoID(video_id int64, user *entity.User) error {
	publication := &entity.Publication{
		VideoID: video_id,
	}
	mysqlDB.AutoMigrate(&entity.Publication{})
	err := mysqlDB.Where("video_id=?", publication.VideoID).First(publication).Error
	if err != nil {
		log.Panicln(err)
		return err
	}

	user.UserID = publication.OwnerID
	UserGetByUID(user)

	return err
}
