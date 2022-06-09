package service

import (
	"log"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/utils"
	"strconv"
	"time"
)

type Comment struct {
	ID         int64
	Content    string
	UserDTO    UserDTO
	CreateDate string
}

func CommentAdd(video_id string, comment_text string, user_id int64) (*Comment, error) {

	var comment_id int64
	var responseComment *Comment
	var user *entity.User
	videoId, _ := strconv.ParseInt(video_id, 10, 64)

	snowFlake, err := utils.NewSnowFlake(0, 0)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	comment_id, err = snowFlake.NextId()
	if err != nil {
		log.Panicln(err)
		return nil, err
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
		log.Println(err)
	}

	dao.UserGetByVideoID(videoId, user)

	userDTO, err1 := UserInfo(user.UserID, user_id)
	if err1 != nil {
		log.Panicln(err)
		return nil, err
	}

	responseComment = &Comment{
		ID:         comment_id,
		Content:    comment_text,
		UserDTO:    *userDTO,
		CreateDate: time.Now().Format("mm-dd"),
	}

	return responseComment, err
}

func CommentDelete(video_id string, comment_id string, user_id int64) (*Comment, error) {

	var comment *entity.Comment
	var responseComment *Comment
	videoId, _ := strconv.ParseInt(video_id, 10, 64)
	commentId, _ := strconv.ParseInt(comment_id, 10, 64)
	var user *entity.User

	comment.CommentID = commentId

	err := dao.CommentGetByCommentID(comment)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = dao.CommentDelete(comment)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dao.UserGetByVideoID(videoId, user)

	userDTO, err1 := UserInfo(user.UserID, user_id)
	if err1 != nil {
		log.Panicln(err1)
		return nil, err1
	}

	responseComment = &Comment{
		ID:         commentId,
		Content:    comment.Content,
		UserDTO:    *userDTO,
		CreateDate: comment.CreateTime.Format("mm-dd"),
	}

	return responseComment, err

}

func CommentList(user_id int64, video_id string) ([]Comment, error) {

	var user *entity.User
	videoId, _ := strconv.ParseInt(video_id, 10, 64)
	var userDTO *UserDTO
	var commentList []entity.Comment
	var commentVoList []Comment

	dao.UserGetByVideoID(videoId, user)

	userDTO, err := UserInfo(user.UserID, user_id)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	err1 := dao.CommentListGetByVideoId(videoId, &commentList)

	if err1 != nil {
		log.Println(err)
		return nil, err1
	}

	for i := 0; i < len(commentList); i++ {
		commentVoList[i].ID = commentList[i].CommentID
		commentVoList[i].Content = commentList[i].Content
		commentVoList[i].UserDTO = *userDTO
		commentVoList[i].CreateDate = commentList[i].CreateTime.Format("mm-dd")
	}

	return commentVoList, err1

}
