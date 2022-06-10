package entity

import "time"

type Comment struct {
	ID        int64  `gorm:"id"`
	CommentID int64  `gorm:"Column:comment_id"`
	UserID    int64  `gorm:"Column:user_id"`
	VideoID   int64  `gorm:"Column:video_id"`
	Content   string `gorm:"Column:content"`
	Status    int16  `gorm:"Column:status"`
	CreatedAt *time.Time
	DeletedAt *time.Time
}
