package entity

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        int64  `gorm:"Column:id;primary_key"`
	CommentID int64  `gorm:"Column:comment_id"`
	UserID    int64  `gorm:"Column:user_id"`
	VideoID   int64  `gorm:"Column:video_id"`
	Content   string `gorm:"Column:content"`
}
