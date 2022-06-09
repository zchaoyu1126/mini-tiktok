package entity

import "time"

type Comment struct {
	CommentID  int64      `gorm:"Column:comment_id"`
	UserID     int64      `gorm:"Column:user_id"`
	VideoID    int64      `gorm:"Column:video_id"`
	Content    string     `gorm:"Column:content"`
	CreateTime *time.Time `gorm:"Column:create_time"`
	Status     int16      `gorm:"Column:status"`
}
